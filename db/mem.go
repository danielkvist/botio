package db

import (
	"fmt"

	"github.com/danielkvist/botio/proto"
)

// Mem is a mocked-up database for testing.
type Mem map[string]string

// Connect simulates a connection with a database.
func (m Mem) Connect() error {
	return nil
}

// Add receives a *proto.BotCommand and adds it
// to the map using the Command as key and the
// Response as a value.
func (m Mem) Add(cmd *proto.BotCommand) error {
	el := cmd.GetCmd().GetCommand()
	val := cmd.GetResp().GetResponse()
	m[el] = val
	return nil
}

// Get receives a *proto.Command and returns if exists
// the respective *proto.BotCommand.
func (m Mem) Get(cmd *proto.Command) (*proto.BotCommand, error) {
	el := cmd.GetCommand()
	val, ok := m[el]
	if !ok {
		return nil, fmt.Errorf("command %q not found", el)
	}

	return &proto.BotCommand{
		Cmd: &proto.Command{
			Command: el,
		},
		Resp: &proto.Response{
			Response: val,
		},
	}, nil
}

// GetAll ranges over the map and returns a *proto.BotCommands
// with all the *proto.BotCommand found.
func (m Mem) GetAll() (*proto.BotCommands, error) {
	var commands []*proto.BotCommand

	for k, v := range m {
		c := &proto.BotCommand{
			Cmd: &proto.Command{
				Command: k,
			},
			Resp: &proto.Response{
				Response: v,
			},
		}

		commands = append(commands, c)
	}

	return &proto.BotCommands{
		Commands: commands,
	}, nil
}

// Remove removes a *proto.BotCommand from the map.
func (m Mem) Remove(cmd *proto.Command) error {
	delete(m, cmd.GetCommand())
	return nil
}

// Update updates the Response of an existing *proto.BotCommand
// with the Response of the received *proto.BotCommand.
// If the *proto.BotCommand didn't exists it adds it.
func (m Mem) Update(cmd *proto.BotCommand) error {
	m[cmd.GetCmd().GetCommand()] = cmd.GetResp().GetResponse()
	return nil
}

// Close deletes all the keys from the map.
func (m Mem) Close() error {
	for k := range m {
		delete(m, k)
	}

	return nil
}
