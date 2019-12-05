package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// Client represents a gRPC BotioCLient with basic CRUD methods.
type Client interface {
	AddCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	GetCommand(context.Context, *proto.Command) (*proto.BotCommand, error)
	ListCommands(context.Context, *empty.Empty) (*proto.BotCommands, error)
	UpdateCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	DeleteCommand(context.Context, *proto.Command) (*empty.Empty, error)
}

type client struct {
	addr   string
	jwt    string
	conn   *grpc.ClientConn
	client proto.BotioClient
}

// ConnOption represents a connection option for a new Client.
type ConnOption func() (*grpc.ClientConn, error)

// WithInsecureConn receives an URL and creates a new non-secured
// connection that it will be assigned to a new Client at the moment of its creation.
// If something goes wrong while creating the insecure connection it returns
// a non-nil error.
func WithInsecureConn(url string) ConnOption {
	return func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, errors.Wrap(err, "while creating a new insecure grpc.ClientConn")
		}

		return conn, nil
	}
}

// WithTLSSecureConn receives an URL and SSL related filenames and
// creates a TLS secured connection that it will be assigned to a new Client at the
// moment of its creation. If something goes wrong while loading the TLS configuration
// or while reading the received certificates it will return a non-nil error.
func WithTLSSecureConn(url, server, crt, key, ca string) ConnOption {
	return func() (*grpc.ClientConn, error) {
		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, errors.Wrap(err, "while loading client SSL key pair")
		}

		certPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, errors.Wrap(err, "while reading CA certificate")
		}

		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New("faile to append CA certificates")
		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   server,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		})

		conn, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, errors.Wrapf(err, "while creating a new Dial for %q", url)
		}

		return conn, nil
	}
}

// New receives and addres, a JWT token for authentication and an option to the new
// Client's connection that is returned. If something goes wrong while
// applying the connection option it returns a non-nil error.
func New(addr string, jwt string, connOpt ConnOption) (Client, error) {
	c := &client{}
	c.addr = addr

	conn, err := connOpt()
	if err != nil {
		return nil, errors.Wrap(err, "while creating new grpc.ClientConn")
	}

	c.jwt = jwt
	c.conn = conn
	c.client = proto.NewBotioClient(c.conn)
	return c, nil
}

func (c *client) AddCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	command := cmd.GetCmd().GetCommand()
	response := cmd.GetResp().GetResponse()
	if command == "" || response == "" {
		return &empty.Empty{}, errors.New("received BotCommand is invalid")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "token", c.jwt)
	if _, err := c.client.AddCommand(ctx, cmd); err != nil {
		return &empty.Empty{}, errors.Wrapf(err, "while adding BotCommand")
	}

	return &empty.Empty{}, nil
}

func (c *client) GetCommand(ctx context.Context, cmd *proto.Command) (*proto.BotCommand, error) {
	command := cmd.GetCommand()
	if command == "" {
		return nil, errors.New("received an empty Command")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "token", c.jwt)
	return c.client.GetCommand(ctx, cmd)
}

func (c *client) ListCommands(ctx context.Context, _ *empty.Empty) (*proto.BotCommands, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "token", c.jwt)
	return c.client.ListCommands(ctx, &empty.Empty{})
}

func (c *client) UpdateCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	command := cmd.GetCmd().GetCommand()
	response := cmd.GetResp().GetResponse()
	if command == "" || response == "" {
		return &empty.Empty{}, errors.New("received BotCommand is invalid")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "token", c.jwt)
	return c.client.UpdateCommand(ctx, cmd)
}

func (c *client) DeleteCommand(ctx context.Context, cmd *proto.Command) (*empty.Empty, error) {
	command := cmd.GetCommand()
	if command == "" {
		return &empty.Empty{}, errors.New("received BotCommand is invalid")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "token", c.jwt)
	return c.client.DeleteCommand(ctx, cmd)
}
