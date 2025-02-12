package rpc

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/SlimeMutation/rpc-service/protobuf/wallet"
)

const MaxRecvMessageSize = 1024 * 1024 * 300

type RpcServerConfig struct {
	GrpcHostname string
	GrpcPort     int
}

type RpcServer struct {
	*RpcServerConfig
	wallet.UnimplementedWalletServiceServer
	stopped atomic.Bool
}

func (s *RpcServer) Stop(ctx context.Context) error {
	s.stopped.Store(true)
	return nil
}

func (s *RpcServer) Stopped() bool {
	return s.stopped.Load()
}

func NewRpcServer(config *RpcServerConfig) (*RpcServer, error) {
	return &RpcServer{
		RpcServerConfig: config,
	}, nil
}

func (s *RpcServer) Start() error {
	go func(s *RpcServer) {
		addr := fmt.Sprintf("%s:%d", s.GrpcHostname, s.GrpcPort)
		fmt.Println("start rpc services", "addr", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("Could not start tcp listener. ")
		}

		opt := grpc.MaxRecvMsgSize(MaxRecvMessageSize)

		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				nil,
			),
		)
		reflection.Register(gs)

		wallet.RegisterWalletServiceServer(gs, s)

		fmt.Println("Grpc info", "port", s.GrpcPort, "address", listener.Addr())
		if err := gs.Serve(listener); err != nil {
			fmt.Println("Could not GRPC services")
		}
	}(s)
	return nil
}
