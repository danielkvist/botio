package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client interface {
	AddCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	GetCommand(context.Context, *proto.Command) (*proto.BotCommand, error)
	ListCommands(context.Context, *empty.Empty) (*proto.BotCommands, error)
	UpdateCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	DeleteCommand(context.Context, *proto.Command) (*empty.Empty, error)
}

type client struct {
	addr   string
	conn   *grpc.ClientConn
	client proto.BotioClient
}

type ConnOption func() (*grpc.ClientConn, error)

func WithInsecureConn(url string) ConnOption {
	return func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("while creating a new insecure grpc.ClientConn: %v", err)
		}

		return conn, nil
	}
}

func WithTLSSecureConn(url, server, crt, key, ca string) ConnOption {
	return func() (*grpc.ClientConn, error) {
		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, fmt.Errorf("while loading client SSL key pair: %v", err)
		}

		certPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("while reading CA certificate: %v", err)
		}

		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("faile to append CA certificates")
		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   server,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		})

		conn, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("while creating a new Dial for %q: %v", url, err)
		}

		return conn, nil
	}
}

func New(addr string, connOpt ConnOption) (Client, error) {
	c := &client{}
	c.addr = addr

	conn, err := connOpt()
	if err != nil {
		return nil, fmt.Errorf("while creating new grpc.ClientConn: %v", err)
	}

	c.conn = conn
	c.client = proto.NewBotioClient(c.conn)
	return c, nil
}

func (c *client) AddCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	command := cmd.GetCmd().GetCommand()
	response := cmd.GetResp().GetResponse()
	if command == "" || response == "" {
		return &empty.Empty{}, fmt.Errorf("received proto.BotCommand to add has an invalid command=%q or response=%q", command, response)
	}

	if _, err := c.client.AddCommand(ctx, cmd); err != nil {
		return &empty.Empty{}, fmt.Errorf("while adding proto.BotCommand: %v", err)
	}

	return &empty.Empty{}, nil
}

func (c *client) GetCommand(ctx context.Context, cmd *proto.Command) (*proto.BotCommand, error) {
	command := cmd.GetCommand()
	if command == "" {
		return nil, fmt.Errorf("received an proto.Command with no command")
	}

	return c.client.GetCommand(ctx, cmd)
}

func (c *client) ListCommands(ctx context.Context, _ *empty.Empty) (*proto.BotCommands, error) {
	return c.client.ListCommands(ctx, &empty.Empty{})
}

func (c *client) UpdateCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	command := cmd.GetCmd().GetCommand()
	response := cmd.GetResp().GetResponse()
	if command == "" || response == "" {
		return nil, fmt.Errorf("received proto.BotCommand to update has an invalid command=%q or response=%q", command, response)
	}

	return c.client.UpdateCommand(ctx, cmd)
}

func (c *client) DeleteCommand(ctx context.Context, cmd *proto.Command) (*empty.Empty, error) {
	command := cmd.GetCommand()
	if command == "" {
		return nil, fmt.Errorf("received an empty proto.Command")
	}

	return c.client.DeleteCommand(ctx, cmd)
}
