package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"sshx/internal/handler"
	"sshx/internal/svc"
	"sshx/internal/types"
)

var configFile = flag.String("f", "etc/conf.yaml", "config file path")

func main() {

	flag.Parse()

	var c types.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	log.Printf("config: %v \n", svcCtx.Config)

	err := handler.SyncExecHandler(svcCtx)
	if err != nil {
		log.Fatal(err)
	}

}
