package main

import (
	"flag"
	"fmt"

	"github.com/SlimeMutation/rpc-service/config"
	services "github.com/SlimeMutation/rpc-service/server"
)

func main() {
	fmt.Println("running grpc services...")
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	cfg, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	grpcServerCfg := &services.RpcServerConfig{
		GrpcHostname: cfg.Server.Host,
		GrpcPort:     cfg.Server.Port,
	}
	rpcServer, err := services.NewRpcServer(grpcServerCfg)
	if err != nil {
		panic(err)
	}
	rpcServer.Start()

	<-(chan struct{})(nil)
}
