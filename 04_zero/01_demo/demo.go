package main

import (
	"flag"
	"fmt"

	"base_go/04_zero/01_demo/internal/config"
	"base_go/04_zero/01_demo/internal/handler"
	"base_go/04_zero/01_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "/Users/thor/GolandProjects/base_go/04_zero/01_demo/etc/demo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
