package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/danielkvist/botio/proto"
	"github.com/danielkvist/botio/server"

	"github.com/golang/protobuf/ptypes/empty"
)

func TestAddCommand(t *testing.T) {
	closeCh := make(chan struct{})
	c := testClient(t, closeCh)

	defer func() {
		close(closeCh)
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
	closeCh := make(chan struct{})
	c := testClient(t, closeCh)

	defer func() {
		close(closeCh)
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
	closeCh := make(chan struct{})
	c := testClient(t, closeCh)

	defer func() {
		close(closeCh)
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
	closeCh := make(chan struct{})
	c := testClient(t, closeCh)

	defer func() {
		close(closeCh)
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
	closeCh := make(chan struct{})
	c := testClient(t, closeCh)

	defer func() {
		close(closeCh)
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

func testClient(t *testing.T, closeCh <-chan struct{}) Client {
	t.Helper()

	s, err := server.New(
		server.WithTestDB(),
		server.WithRistrettoCache(262144000),
		server.WithListener(":33333"),
		server.WithInsecureGRPCServer(),
		server.WithTextLogger(&bytes.Buffer{}),
	)
	if err != nil {
		t.Fatalf("while creating a new Server for testing: %v", err)
	}

	go func(t *testing.T) {
		s.Serve()
	}(t)

	go func() {
		<-closeCh
		s.CloseList()
	}()

	c, err := New(":33333", WithInsecureConn(":33333"))
	return c
}
