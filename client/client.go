package client

import (
	"context"
	"fmt"

	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
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

func New(addr string, conn *grpc.ClientConn) Client {
	c := &client{}

	c.addr = addr
	c.conn = conn
	c.client = proto.NewBotioClient(c.conn)
	return c
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
		return nil, fmt.Errorf("received an proto.Command with no command")
	}

	return c.client.DeleteCommand(ctx, cmd)
}
