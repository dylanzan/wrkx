package sshx

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net"
	"os"
	"time"
)

type Cli struct {
	IP       string
	UserName string
	Password string
	Port     int
	client   *ssh.Client
}

func NewSSHX(ip, username, password string, port ...int) (*Cli, error) {
	cli := &Cli{
		IP:       ip,
		UserName: username,
		Password: password,
	}
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}

	return cli, nil
}

// connect to ssh server
func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)

	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = client
	//log.Printf("%v is connected.", addr)
	return nil
}

func (c *Cli) Run(shell string) (string, error) {
	session, err := c.newSession()

	if err != nil {
		return "", err
	}
	defer session.Close()

	buf, err := session.CombinedOutput(shell)

	return string(buf), err
}

func (c *Cli) RunTerminal(shell string) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}

// RunTerminalSession run terminal session
func (c *Cli) runTerminalSession(session *ssh.Session, shell string) error {
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	terminal, termHeight, err := terminal.GetSize(fd)

	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", termHeight, terminal, modes); err != nil {
		return err
	}

	session.Run(shell)
	return nil
}

func (c *Cli) newSession() (*ssh.Session, error) {
	if c.client == nil {
		err := c.connect()
		if err != nil {
			return nil, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *Cli) EnterTerminal() error {
	session, err := c.newSession()
	if err != nil {
		return err
	}

	defer session.Close()
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	terminal, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", termHeight, terminal, modes); err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	err = session.Wait()
	return err
}

func (c *Cli) Enter(w io.Writer, r io.Reader) error {
	session, err := c.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	fd := int(os.Stdin.Fd())

	session.Stdout = w
	session.Stderr = os.Stdin
	session.Stdin = r

	terminal, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm-256color", termHeight, terminal, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	err = session.Wait()

	return err

}
