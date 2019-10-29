package server

import (
	"context"
	"fmt"

	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddCommand tries to add a received command to the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) AddCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if err := s.db.Add(cmd); err != nil {
			return &empty.Empty{}, status.Error(codes.Internal, fmt.Sprintf("while adding command %q: %v", cmd.GetCmd().GetCommand(), err))
		}
	}

	return &empty.Empty{}, nil
}

// GetCommand tries to get the specified command from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) GetCommand(ctx context.Context, cmd *proto.Command) (*proto.BotCommand, error) {
	var c *proto.BotCommand
	var err error

	select {
	case <-ctx.Done():
		return &proto.BotCommand{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		c, err = s.db.Get(cmd)
		if err != nil {
			return &proto.BotCommand{}, status.Error(codes.Internal, fmt.Sprintf("while getting command %q: %v", cmd.GetCommand(), err))
		}
	}

	return c, nil
}

// ListCommands tries to get all the commands from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) ListCommands(ctx context.Context, _ *empty.Empty) (*proto.BotCommands, error) {
	var commands *proto.BotCommands
	var err error

	select {
	case <-ctx.Done():
		return &proto.BotCommands{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		commands, err = s.db.GetAll()
		if err != nil {
			return &proto.BotCommands{}, status.Error(codes.Internal, fmt.Sprintf("while getting commands: %v", err))
		}
	}

	return commands, nil
}

// UpdateCommand tries to update the specified command to the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) UpdateCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if err := s.db.Update(cmd); err != nil {
			return &empty.Empty{}, status.Error(codes.Internal, fmt.Sprintf("while updating command %q: %v", cmd.GetCmd().GetCommand(), err))
		}
	}

	return &empty.Empty{}, nil
}

// DeleteCommand tries to remove the specified command from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) DeleteCommand(ctx context.Context, cmd *proto.Command) (*empty.Empty, error) {
	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if err := s.db.Remove(cmd); err != nil {
			return &empty.Empty{}, status.Error(codes.Internal, fmt.Sprintf("while deleting command %q: %v", cmd.GetCommand(), err))
		}
	}

	return &empty.Empty{}, nil
}
