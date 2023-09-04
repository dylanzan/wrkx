package sshx

import (
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"sshx/internal/types"
	"strings"
	"testing"
)

func TestSSHX(t *testing.T) {

	configFile := "/Users/dylan.zan/Documents/goproject/sshx/etc/conf.yaml"

	var c types.Config
	conf.MustLoad(configFile, &c)

	serverinfos := strings.Split(c.Servers[0], ":")

	cli, err := NewSSHX(serverinfos[0], serverinfos[2], serverinfos[3])

	if err != nil {
		t.Error(err)
	}

	stdoutstr, err := cli.Run("curl -L ip.iizone.com.cn")
	if err != nil {
		t.Error(err)
	}
	log.Println(stdoutstr)

}
