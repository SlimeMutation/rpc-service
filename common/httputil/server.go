package httputil

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"sync/atomic"
)

type HTTPServer struct {
	listener net.Listener
	srv      *http.Server
	closed   atomic.Bool
}

type HTTPOption func(srv *HTTPServer) error

func StartHttpServer(addr string, handler http.Handler, opts ...HTTPOption) (*HTTPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Init listener failed: ", err)
		return nil, errors.New("Init listener failed: " + err.Error())
	}
	srvCtx, srvCancel := context.WithCancel(context.Background())
	srv := &http.Server{
		Handler:           handler,
		ReadTimeout:       timeouts.ReadTimeout,
		ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
		WriteTimeout:      timeouts.WriteTimeout,
		IdleTimeout:       timeouts.IdleTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return srvCtx
		},
	}
	out := &HTTPServer{
		listener: listener,
		srv:      srv,
	}

	for _, opt := range opts {
		if err := opt(out); err != nil {
			srvCancel()
			fmt.Println("One of http op failed: ", err)
			return nil, errors.New("One of http op failed: " + err.Error())
		}
	}
	go func() {
		err := out.srv.Serve(out.listener)
		srvCancel()
		if errors.Is(err, http.ErrServerClosed) {
			out.closed.Store(true)
		} else {
			fmt.Println("unknown error: ", err)
			panic("unknown error")
		}
	}()
	return out, nil
}

func (s *HTTPServer) Closed() bool {
	return s.closed.Load()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		if errors.Is(err, ctx.Err()) {
			return s.Close()
		}
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *HTTPServer) Close() error {
	return s.srv.Close()
}

func (s *HTTPServer) Addr() net.Addr {
	return s.listener.Addr()
}

func WithMaxHeaderBytes(maxHeaderBytes int) HTTPOption {
	return func(s *HTTPServer) error {
		s.srv.MaxHeaderBytes = maxHeaderBytes
		return nil
	}
}
