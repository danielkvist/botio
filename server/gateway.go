package server

import (
	"context"
	"net/http"

	"github.com/danielkvist/botio/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func (s *server) jsonGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := proto.RegisterBotioHandlerFromEndpoint(ctx, mux, s.listener.Addr().String(), options)
	if err != nil {
		return errors.Wrapf(err, "while registering Botio HTTP handler")
	}

	return http.ListenAndServe(s.httpPort, mux)
}
