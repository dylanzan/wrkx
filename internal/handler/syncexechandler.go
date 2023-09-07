package handler

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"log"
	"os"
	"sshx/internal/cores/sshx"
	"sshx/internal/svc"
	"sshx/internal/tools/ostools"
	"strings"
	"sync"
	"time"
)

var (
	_wg      = &sync.WaitGroup{}
	_poolMap = make(map[string]*sync.Pool)
)

func SyncExecHandler(svcCtx *svc.ServiceContext, command string) error {
	if len(svcCtx.Config.Servers) == 0 {
		return errors.New("no server")
	}

	if svcCtx.Config.IsCheckOs {
		serverVersions := make(map[string]struct{})
		for _, serverInfo := range svcCtx.Config.Servers {
			host, port, user, password := _getServerInfo(serverInfo)
			cli, err := sshx.NewSSHX(host, user, password, port)
			if err != nil {
				return err
			}
			ovstr, err := ostools.GetOSVersionName(cli)
			log.Printf("%s:%d %s", host, port, ovstr)
			if err != nil {
				return err
			}
			serverVersions[ovstr] = struct{}{}
		}
		if len(serverVersions) > 1 {
			return errors.New("different os version")
		}
	}

	for _, serverInfo := range svcCtx.Config.Servers {

		host, port, user, password := _getServerInfo(serverInfo)
		cli, err := sshx.NewSSHX(host, user, password, port)
		if err != nil {
			return err
		}
		keyStr := fmt.Sprintf("%s:%d", host, port)
		_poolMap[keyStr] = &sync.Pool{
			New: func() interface{} {
				return cli
			},
		}
	}

	//for {
	//cmdStr := strings.TrimRight(_commandScan(), "\n")
	//if cmdStr == "exit" {
	//	os.Exit(0)
	//}
	for serverKey := range _poolMap {
		_wg.Add(1)
		go _execCommand(serverKey, command, _wg)
	}
	_wg.Wait()
	//}

	return nil
}

// _getServerInfo 获取服务器信息
func _getServerInfo(serverInfo string) (host string, port int, user string, password string) {
	serverInfos := strings.Split(serverInfo, ":")
	return serverInfos[0], cast.ToInt(serverInfos[1]), serverInfos[2], serverInfos[3]
}

// _commandScan 命令行输入
func _commandScan() (input string) {
	fmt.Print(">> ")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if input == "\n" {
		return _commandScan()
	}
	if err != nil {
		fmt.Println("There were errors reading, exiting program.")
		return
	}

	return input
}

// _execCommand 执行命令
func _execCommand(serverKey string, cmdStr string, wg *sync.WaitGroup) (err error) {

	defer wg.Done()
	cli := _poolMap[serverKey].Get().(*sshx.Cli)
	defer _poolMap[serverKey].Put(cli)

	stdout, err := cli.Run(cmdStr)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	fmt.Printf("----------%s----------\n%s\n----------%s done.----------\n", serverKey, stdout, serverKey)
	return nil
}
