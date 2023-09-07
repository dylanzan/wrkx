package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"sshx/internal/handler"
	"sshx/internal/svc"
	"sshx/internal/types"
)

var configFile = flag.String("f", "etc/conf-prod.yaml", "config file path")
var command = flag.String("c", "", "command")

func main() {

	flag.Parse()

	var config types.Config
	conf.MustLoad(*configFile, &config)

	svcCtx := svc.NewServiceContext(config)
	log.Printf("config: %v \n", svcCtx.Config)

	err := handler.SyncExecHandler(svcCtx, *command)
	if err != nil {
		log.Fatal(err)
	}
}
