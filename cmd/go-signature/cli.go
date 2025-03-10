package main

import (
	"context"
	"fmt"
	"github.com/SlimeMutation/rpc-service/services/rest"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"

	"github.com/SlimeMutation/rpc-service/common/cliapp"
	"github.com/SlimeMutation/rpc-service/common/opio"
	"github.com/SlimeMutation/rpc-service/config"
	"github.com/SlimeMutation/rpc-service/database"
	flags2 "github.com/SlimeMutation/rpc-service/flags"
	"github.com/SlimeMutation/rpc-service/services/rpc"
)

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	fmt.Println("running grpc services...")
	cfg := config.NewConfig(ctx)
	grpcServerCfg := &rpc.RpcServerConfig{
		GrpcHostname: cfg.RpcServer.Host,
		GrpcPort:     cfg.RpcServer.Port,
	}
	db, err := database.NewDB(ctx.Context, cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}
	return rpc.NewRpcServer(db, grpcServerCfg)
}

func runMigrations(ctx *cli.Context) error {
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	log.Info("running migrations...")
	cfg := config.NewConfig(ctx)
	db, err := database.NewDB(ctx.Context, cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Error("fail to close database", "err", err)
		}
	}(db)
	return db.ExecuteSQLMigration(cfg.Migrations)
}

func runRestApi(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("running api...")
	cfg := config.NewConfig(ctx)
	return rest.NewApi(ctx.Context, &cfg)
}

func NewCli(GitCommit string, GitData string) *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              params.VersionWithCommit(GitCommit, GitData),
		Description:          "An exchange wallet scanner services with rpc and rest api services",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "api",
				Flags:       flags,
				Description: "Run api services",
				Action:      cliapp.LifecycleCmd(runRestApi),
			},
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run rpc services",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "migrate",
				Flags:       flags,
				Description: "Run database migrations",
				Action:      runMigrations,
			},
			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
