package db

import (
	"testing"

	"github.com/danielkvist/botio/proto"
)

func TestConnect(t *testing.T) {
	var m Mem
	m = make(map[string]string)
	if err := m.Connect(); err != nil {
		t.Fatalf("while connecting Mem should never fail: %v", err)
	}
}

func TestAdd(t *testing.T) {
	tt := []struct {
		cmd  *proto.Command
		resp *proto.Response
	}{
		{
			cmd: &proto.Command{
				Command: "test",
			},
			resp: &proto.Response{
				Response: "this is a test",
			},
		},
	}

	var m Mem
	m = make(map[string]string, 1)

	for _, tc := range tt {
		command := &proto.BotCommand{
			Cmd:  tc.cmd,
			Resp: tc.resp,
		}
		if err := m.Add(command); err != nil {
			t.Fatalf("while adding command %q: %v", tc.cmd.GetCommand(), err)
		}
	}
}

func TestGet(t *testing.T) {
	commandOne := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "a",
		},
		Resp: &proto.Response{
			Response: "abc",
		},
	}

	commandTwo := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "z",
		},
		Resp: &proto.Response{
			Response: "xyz",
		},
	}

	tt := []struct {
		command          *proto.BotCommand
		expectedCmd      *proto.Command
		expectedResponse *proto.Response
	}{
		{
			command:          commandOne,
			expectedCmd:      commandOne.Cmd,
			expectedResponse: commandOne.Resp,
		},
		{
			command:          commandTwo,
			expectedCmd:      commandTwo.Cmd,
			expectedResponse: commandTwo.Resp,
		},
	}

	var m Mem
	m = map[string]string{
		commandOne.Cmd.Command: commandOne.Resp.Response,
		commandTwo.Cmd.Command: commandTwo.Resp.Response,
	}

	for _, tc := range tt {
		cmd, err := m.Get(tc.command.Cmd)
		if err != nil {
			t.Fatalf("error while getting command %q: %v", tc.command.GetCmd().GetCommand(), err)
		}

		if cmd.GetCmd().GetCommand() != tc.expectedCmd.GetCommand() {
			t.Fatalf("expected command to be %q. got=%v", tc.expectedCmd.GetCommand(), cmd.GetCmd().GetCommand())
		}

		if cmd.GetResp().GetResponse() != tc.expectedResponse.GetResponse() {
			t.Fatalf("expected command response to be %q. got=%v", tc.expectedResponse.GetResponse(), cmd.GetResp().GetResponse())
		}
	}
}

func TestGetAll(t *testing.T) {
	commandOne := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "a",
		},
		Resp: &proto.Response{
			Response: "abc",
		},
	}

	commandTwo := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "z",
		},
		Resp: &proto.Response{
			Response: "xyz",
		},
	}

	var m Mem
	m = map[string]string{
		commandOne.Cmd.Command: commandOne.Resp.Response,
		commandTwo.Cmd.Command: commandTwo.Resp.Response,
	}

	commands, err := m.GetAll()
	if err != nil {
		t.Fatalf("while getting all the commands: %v", err)
	}

	if len(commands.GetCommands()) != 2 {
		t.Fatalf("expected to get 2 commands. got=%v", len(commands.GetCommands()))
	}
}

func TestRemove(t *testing.T) {
	command := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "a",
		},
		Resp: &proto.Response{
			Response: "abc",
		},
	}

	var m Mem
	m = map[string]string{
		command.Cmd.Command: command.Resp.Response,
	}

	if len(m) != 1 {
		t.Fatalf("expected map to have 1 item. got=%v", len(m))
	}

	if err := m.Remove(command.GetCmd()); err != nil {
		t.Fatalf("while removing command %q: %v", command.GetCmd().GetCommand(), err)
	}

	if len(m) != 0 {
		t.Fatalf("expected map to have 0 item. got=%v", len(m))
	}
}

func TestUpdate(t *testing.T) {
	oldCommand := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "a",
		},
		Resp: &proto.Response{
			Response: "abc",
		},
	}

	newCommand := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "a",
		},
		Resp: &proto.Response{
			Response: "xyz",
		},
	}

	var m Mem
	m = map[string]string{
		oldCommand.Cmd.Command: oldCommand.Resp.Response,
	}

	if err := m.Update(newCommand); err != nil {
		t.Fatalf("while updating command responde: %v", err)
	}

	response, ok := m[oldCommand.Cmd.GetCommand()]
	if !ok {
		t.Fatalf("while checking if command %q exists. command not found", oldCommand.Cmd.GetCommand())
	}

	if response != newCommand.GetResp().GetResponse() {
		t.Fatalf("expected command to have update response %q. got=%q", newCommand.GetResp().GetResponse(), response)
	}
}

func TestClose(t *testing.T) {
	var m Mem
	m = make(map[string]string)
	m.Add(&proto.BotCommand{
		Cmd: &proto.Command{
			Command: "Hi",
		},
		Resp: &proto.Response{
			Response: "Hello, World!",
		},
	})

	if err := m.Close(); err != nil {
		t.Fatalf("while closing Mem should never fail: %v", err)
	}
}
