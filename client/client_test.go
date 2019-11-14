package client

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/danielkvist/botio/proto"
	"github.com/danielkvist/botio/server"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestAddCommand(t *testing.T) {
	close := make(chan struct{})
	c := testClient(t, close)

	defer func() {
		close <- struct{}{}
	}()

	tt := []struct {
		name           string
		command        *proto.BotCommand
		expectedToFail bool
	}{
		{
			name:           "empty command",
			command:        &proto.BotCommand{},
			expectedToFail: true,
		},
		{
			name: "command with both fields",
			command: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "start",
				},
				Resp: &proto.Response{
					Response: "Hi",
				},
			},
		},
		{
			name: "command without command",
			command: &proto.BotCommand{
				Resp: &proto.Response{
					Response: "Hi",
				},
			},
			expectedToFail: true,
		},
		{
			name: "command without response",
			command: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "start",
				},
			},
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := c.AddCommand(context.TODO(), tc.command); err != nil {
				if tc.expectedToFail {
					t.Skipf("while adding proto.BotCommand failed as expected: %v", err)
				}

				t.Fatalf("while adding proto.BotCommand: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("test expected to fail did not fail")
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	close := make(chan struct{})
	c := testClient(t, close)

	defer func() {
		close <- struct{}{}
	}()

	command := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	if _, err := c.AddCommand(context.TODO(), command); err != nil {
		t.Fatalf("while adding command for testing: %v", err)
	}

	cmd, err := c.GetCommand(context.TODO(), command.GetCmd())
	if err != nil {
		t.Fatalf("while getting command previously added: %v", err)
	}

	if cmd.GetCmd().GetCommand() != command.GetCmd().GetCommand() {
		t.Fatalf("expected to get command %q. got=%q command", command.GetCmd().GetCommand(), cmd.GetCmd().GetCommand())
	}

	if cmd.GetResp().GetResponse() != command.GetResp().GetResponse() {
		t.Fatalf("expected to get command %q with response %q. got=%q command with response=%q", command.GetCmd().GetCommand(), command.GetResp().GetResponse(), cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
	}

	if _, err := c.GetCommand(context.TODO(), &proto.Command{Command: ""}); err == nil {
		t.Fatalf("bad command should not be handled")
	}

	if _, err := c.GetCommand(context.TODO(), nil); err == nil {
		t.Fatalf("nil command should not be handled")
	}
}

func TestListCommand(t *testing.T) {
	close := make(chan struct{})
	c := testClient(t, close)

	defer func() {
		close <- struct{}{}
	}()

	commandOne := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	commandTwo := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "close",
		},
		Resp: &proto.Response{
			Response: "goodbye",
		},
	}

	commands := []*proto.BotCommand{commandOne, commandTwo}
	for _, cmd := range commands {
		if _, err := c.AddCommand(context.TODO(), cmd); err != nil {
			t.Fatalf("while adding command %q for testing: %v", cmd.GetCmd().GetCommand(), err)
		}
	}

	list, err := c.ListCommands(context.TODO(), &empty.Empty{})
	if err != nil {
		t.Fatalf("while listening commands: %v", err)
	}

	if len(list.GetCommands()) != 2 {
		t.Fatalf("expected to get a list with %v elements. got=%v", 2, len(list.GetCommands()))
	}
}

func TestUpdateCommand(t *testing.T) {
	close := make(chan struct{})
	c := testClient(t, close)

	defer func() {
		close <- struct{}{}
	}()

	command := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	newCommand := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "goodbye",
		},
	}

	if _, err := c.AddCommand(context.TODO(), command); err != nil {
		t.Fatalf("while adding command for testing: %v", err)
	}

	cmd, err := c.GetCommand(context.TODO(), command.GetCmd())
	if err != nil {
		t.Fatalf("while getting command previously added: %v", err)
	}

	if cmd.GetCmd().GetCommand() != command.GetCmd().GetCommand() {
		t.Fatalf("expected to get command %q. got=%q command", command.GetCmd().GetCommand(), cmd.GetCmd().GetCommand())
	}

	if cmd.GetResp().GetResponse() != command.GetResp().GetResponse() {
		t.Fatalf("expected to get command %q with response %q. got=%q command with response=%q", command.GetCmd().GetCommand(), command.GetResp().GetResponse(), cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
	}

	if _, err := c.UpdateCommand(context.TODO(), newCommand); err != nil {
		t.Fatalf("while updating command: %v", err)
	}

	cmd, err = c.GetCommand(context.TODO(), command.GetCmd())
	if err != nil {
		t.Fatalf("while getting command previously updated: %v", err)
	}

	if cmd.GetCmd().GetCommand() != newCommand.GetCmd().GetCommand() {
		t.Fatalf("expected to get command %q. got=%q command", newCommand.GetCmd().GetCommand(), cmd.GetCmd().GetCommand())
	}

	if cmd.GetResp().GetResponse() != newCommand.GetResp().GetResponse() {
		t.Fatalf("expected to get command %q with response %q. got=%q command with response=%q", newCommand.GetCmd().GetCommand(), newCommand.GetResp().GetResponse(), cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
	}

	if _, err := c.GetCommand(context.TODO(), nil); err == nil {
		t.Fatalf("nil command should not be handled")
	}
}

func TestDeleteCommand(t *testing.T) {
	close := make(chan struct{})
	c := testClient(t, close)

	defer func() {
		close <- struct{}{}
	}()

	command := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	if _, err := c.AddCommand(context.TODO(), command); err != nil {
		t.Fatalf("while adding command for testing: %v", err)
	}

	list, err := c.ListCommands(context.TODO(), &empty.Empty{})
	if err != nil {
		t.Fatalf("while listening commands: %v", err)
	}

	if len(list.GetCommands()) != 1 {
		t.Fatalf("expected to get a list with %v elements. got=%v", 1, len(list.GetCommands()))
	}

	if _, err := c.DeleteCommand(context.TODO(), command.GetCmd()); err != nil {
		t.Fatalf("while deleting command %q: %v", command.GetCmd(), err)
	}

	list, err = c.ListCommands(context.TODO(), &empty.Empty{})
	if err != nil {
		t.Fatalf("while listening commands: %v", err)
	}

	if len(list.GetCommands()) != 0 {
		t.Fatalf("expected to get a list with %v elements. got=%v", 0, len(list.GetCommands()))
	}
}

func testClient(t *testing.T, close <-chan struct{}) client {
	t.Helper()

	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)
	s, err := server.New(server.WithTestDB())
	if err != nil {
		t.Fatalf("while creating a new Server for testing: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterBotioServer(srv, s)
	go func(t *testing.T) {
		if err := srv.Serve(listener); err != nil {
			t.Fatalf("while listening with a BotioServer: %v", err)
		}
	}(t)

	dialer := func(_ string, _ time.Duration) (net.Conn, error) { return listener.Dial() }
	conn, err := grpc.DialContext(context.TODO(), "test", grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("while creatin a gRPC dial: %v", err)
	}

	go func() {
		<-close
		defer conn.Close()
	}()

	var c client
	c.conn = conn
	c.addr = listener.Addr().String()
	c.client = proto.NewBotioClient(conn)

	return c
}
