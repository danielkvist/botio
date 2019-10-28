package server

import (
	"context"
	"testing"

	"github.com/danielkvist/botio/proto"
)

func TestAddCommand(t *testing.T) {
	tt := []struct {
		name           string
		command        *proto.BotCommand
		expectedToFail bool
	}{
		{
			name: "without command",
		},
		{
			name: "with command",
			command: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "start",
				},
				Resp: &proto.Response{
					Response: "hi",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := testServer(t)
			_, err := s.AddCommand(context.TODO(), tc.command)
			if err != nil {
				if tc.expectedToFail {
					t.Skipf("add command operation failed as expected: %v", err)
				}

				t.Fatalf("while adding command: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("add command operation not failed as expected")
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	command := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	s := testServer(t)
	if _, err := s.AddCommand(context.TODO(), command); err != nil {
		t.Fatalf("while adding command to the database: %v", err)
	}

	cmd, err := s.GetCommand(context.TODO(), command.GetCmd())
	if err != nil {
		t.Fatalf("while getting command %q: %v", command.GetCmd().GetCommand(), err)
	}

	if cmd.GetCmd().GetCommand() != command.GetCmd().GetCommand() {
		t.Fatalf("expected command %q. got=%q", command.GetCmd().GetCommand(), cmd.GetCmd().GetCommand())
	}

	if cmd.GetResp().GetResponse() != command.GetResp().GetResponse() {
		t.Fatalf("expected command with response %q. got=%q", command.GetResp().GetResponse(), cmd.GetResp().GetResponse())
	}
}

func TestListCommands(t *testing.T) {
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
			Command: "end",
		},
		Resp: &proto.Response{
			Response: "goodbye",
		},
	}

	s := testServer(t)
	if _, err := s.AddCommand(context.TODO(), commandOne); err != nil {
		t.Fatalf("while adding command %q: %v", commandOne.GetCmd().GetCommand(), err)
	}

	if _, err := s.AddCommand(context.TODO(), commandTwo); err != nil {
		t.Fatalf("while adding command %q: %v", commandTwo.GetCmd().GetCommand(), err)
	}

	commands, err := s.ListCommands(context.TODO(), &proto.Void{})
	if err != nil {
		t.Fatalf("while listing commands: %v", err)
	}

	if len(commands.GetCommands()) != 2 {
		t.Fatalf("expected to get a list of commands with %v elements. got=%v elements", 2, commands.GetCommands())
	}

	if commands.GetCommands()[0].GetCmd().GetCommand() != commandOne.GetCmd().GetCommand() {
		t.Fatalf("expected first command to be command %q. got=%q", commands.GetCommands()[0].GetCmd().GetCommand(), commandOne.GetCmd().GetCommand())
	}

	if commands.GetCommands()[1].GetCmd().GetCommand() != commandTwo.GetCmd().GetCommand() {
		t.Fatalf("expected second command to be command %q. got=%q", commands.GetCommands()[1].GetCmd().GetCommand(), commandTwo.GetCmd().GetCommand())
	}
}

func TestUpdateCommand(t *testing.T) {
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
			Response: "hello",
		},
	}

	s := testServer(t)
	if _, err := s.AddCommand(context.TODO(), command); err != nil {
		t.Fatalf("while adding command: %v", err)
	}

	if _, err := s.UpdateCommand(context.TODO(), newCommand); err != nil {
		t.Fatalf("while updating command: %v", err)
	}

	cmd, err := s.GetCommand(context.TODO(), command.GetCmd())
	if err != nil {
		t.Fatalf("while getting comand %q: %v", command.GetCmd(), err)
	}

	if cmd.GetResp().GetResponse() != newCommand.GetResp().GetResponse() {
		t.Fatalf("expected command to have updated response %q. got=%q", newCommand.GetResp().GetResponse(), cmd.GetResp().GetResponse())
	}
}

func TestDeleteCommand(t *testing.T) {
	s := testServer(t)
	if _, err := s.AddCommand(context.TODO(), &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}); err != nil {
		t.Fatalf("while adding command: %v", err)
	}

	if _, err := s.DeleteCommand(context.TODO(), &proto.Command{Command: "start"}); err != nil {
		t.Fatalf("while deleting command: %v", err)
	}

	commands, err := s.ListCommands(context.TODO(), &proto.Void{})
	if err != nil {
		t.Fatalf("while listing commands: %v", err)
	}

	if len(commands.GetCommands()) != 0 {
		t.Fatalf("expected no commands. got=%v commands", len(commands.GetCommands()))
	}
}

func testServer(t *testing.T) Server {
	t.Helper()
	s, err := New(WithTestDB())
	if err != nil {
		t.Fatalf("while creating a new Server for testing: %v", err)
	}

	if err := s.Connect(); err != nil {
		t.Fatalf("while connecting Server for testing to its database: %v", err)
	}

	return s
}
