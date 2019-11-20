package cache

import (
	"testing"
	"time"

	"github.com/danielkvist/botio/proto"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name           string
		capValue       int64
		expectedToFail bool
	}{
		{
			name:     "with capacity",
			capValue: 1 << 30,
		},
		{
			name:           "without capacity",
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		if _, err := New(tc.capValue); err != nil {
			if tc.expectedToFail {
				t.Skipf("test failed as expected: %v", err)
			}

			t.Fatal(err)
		}

		if tc.expectedToFail {
			t.Fatalf("test expected to fail with capacity %v not failed as expected", tc.capValue)
		}
	}
}

func TestAdd(t *testing.T) {
	tt := []struct {
		name           string
		cmd            *proto.BotCommand
		expectedToFail bool
	}{
		{
			name: "valid bot command",
			cmd: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "start",
				},
				Resp: &proto.Response{
					Response: "hi",
				},
			},
		},
		{
			name: "invalid bot command",
			cmd: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "",
				},
				Resp: &proto.Response{
					Response: "hi",
				},
			},
			expectedToFail: true,
		},
		{
			name: "invalid bot response",
			cmd: &proto.BotCommand{
				Cmd: &proto.Command{
					Command: "start",
				},
				Resp: &proto.Response{
					Response: "",
				},
			},
			expectedToFail: true,
		},
	}

	c, err := New(1 << 30)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		if err := c.Add(tc.cmd); err != nil {
			if tc.expectedToFail {
				t.Skipf("test failed as expected: %v", err)
			}

			t.Fatal(err)
		}

		if tc.expectedToFail {
			t.Fatalf("test expected to fail with command %q and response %q not failed as expected", tc.cmd.GetCmd().GetCommand(), tc.cmd.GetResp().GetResponse())
		}
	}
}

func TestGet(t *testing.T) {
	cmd := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	c, err := New(1 << 30)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Add(cmd); err != nil {
		t.Fatalf("while adding command for testing: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	for i := 0; i <= 1000; i++ {
		command, err := c.Get(cmd.GetCmd())
		if err != nil {
			t.Fatalf("(%v) while getting command %q: %v", i, cmd.GetCmd().GetCommand(), err)
		}

		if command.GetCmd().GetCommand() != cmd.GetCmd().GetCommand() {
			t.Fatalf("(%v) expected to get command %q. got=%q", i, cmd.GetCmd().GetCommand(), command.GetCmd().GetCommand())
		}

		if command.GetResp().Response != cmd.GetResp().GetResponse() {
			t.Fatalf("(%v) expected to get command %q with response %q. got response=%q", i, command.GetCmd().GetCommand(), cmd.GetResp().GetResponse(), command.GetResp().GetResponse())
		}
	}
}

func TestRemove(t *testing.T) {
	cmd := &proto.BotCommand{
		Cmd: &proto.Command{
			Command: "start",
		},
		Resp: &proto.Response{
			Response: "hi",
		},
	}

	c, err := New(1 << 30)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Add(cmd); err != nil {
		t.Fatalf("while adding command for testing: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if err := c.Remove(cmd.GetCmd()); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)
	if _, err = c.Get(cmd.GetCmd()); err == nil {
		t.Fatalf("command %q should have triggered an error", cmd.GetCmd().GetCommand())
	}
}
