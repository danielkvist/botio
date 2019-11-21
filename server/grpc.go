package server

import (
	"context"
	"fmt"
	"time"

	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddCommand tries to add a received command to the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) AddCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	start := time.Now()

	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if err := s.db.Add(cmd); err != nil {
			s.logError(
				"db",
				"Add",
				err.Error(),
				fmt.Sprintf("add BotCommand %q: %q failed", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse()),
			)
			return &empty.Empty{}, status.Error(codes.Internal, "error while adding command")
		}
	}

	s.logInfo(
		"server",
		"AddCommand",
		fmt.Sprintf("BotCommand %q: %q added successfully", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse()),
		time.Since(start),
	)
	return &empty.Empty{}, nil
}

// GetCommand tries to get the specified command from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) GetCommand(ctx context.Context, cmd *proto.Command) (*proto.BotCommand, error) {
	var c *proto.BotCommand
	var err error

	start := time.Now()

	select {
	case <-ctx.Done():
		return &proto.BotCommand{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if ok := s.inCache(cmd); !ok {
			c, err = s.db.Get(cmd)
			if err != nil {
				s.logError(
					"db",
					"Get",
					err.Error(),
					fmt.Sprintf("get BotCommand %q failed", cmd.GetCommand()),
				)
				return &proto.BotCommand{}, status.Error(codes.Internal, "error while getting command")
			}

			if err := s.cache.Add(c); err != nil {
				s.logError(
					"cache",
					"Add",
					err.Error(),
					fmt.Sprintf("add BotCommand %q: %q failed", c.GetCmd().GetCommand(), c.GetResp().GetResponse()),
				)
			}
		} else {
			c, err = s.cache.Get(cmd)
			if err != nil {
				s.logError(
					"cache",
					"Get",
					err.Error(),
					fmt.Sprintf("get BotCommand %q failed", cmd.GetCommand()),
				)
				return &proto.BotCommand{}, status.Error(codes.Internal, "error while getting command")
			}
		}
	}

	s.logInfo(
		"server",
		"GetCommand",
		fmt.Sprintf("BotCommand %q gotten successfully", c.GetCmd().GetCommand()),
		time.Since(start),
	)
	return c, nil
}

// ListCommands tries to get all the commands from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) ListCommands(ctx context.Context, _ *empty.Empty) (*proto.BotCommands, error) {
	var commands *proto.BotCommands
	var err error

	start := time.Now()

	select {
	case <-ctx.Done():
		return &proto.BotCommands{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		commands, err = s.db.GetAll()
		if err != nil {
			s.logError(
				"db",
				"GetAll",
				err.Error(),
				"get BotCommands failed",
			)

			return &proto.BotCommands{}, status.Error(codes.Internal, "error while getting commands")
		}
	}

	s.logInfo(
		"server",
		"ListCommands",
		"BotCommands gotten successfully",
		time.Since(start),
	)
	return commands, nil
}

// UpdateCommand tries to update the specified command to the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) UpdateCommand(ctx context.Context, cmd *proto.BotCommand) (*empty.Empty, error) {
	start := time.Now()

	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if ok := s.inCache(cmd.GetCmd()); ok {
			if err := s.cache.Remove(cmd.GetCmd()); err != nil {
				s.logError(
					"cache",
					"Remove",
					err.Error(),
					fmt.Sprintf("remove BotCommand %q failed", cmd.GetCmd().GetCommand()),
				)
			}
		}

		if err := s.db.Update(cmd); err != nil {
			s.logError(
				"db",
				"Update",
				err.Error(),
				fmt.Sprintf("update BotCommand %q: %q failed", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse()),
			)
			return &empty.Empty{}, status.Error(codes.Internal, "error while updating command")
		}
	}

	s.logInfo(
		"server",
		"UpdateCommand",
		fmt.Sprintf("BotCommand %q: %q updated successfully", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse()),
		time.Since(start),
	)
	return &empty.Empty{}, nil
}

// DeleteCommand tries to remove the specified command from the Server's database. It returns a non-nil error
// if something went wrong or if the context was cancelled.
func (s *server) DeleteCommand(ctx context.Context, cmd *proto.Command) (*empty.Empty, error) {
	start := time.Now()

	select {
	case <-ctx.Done():
		return &empty.Empty{}, status.Error(codes.Canceled, ctx.Err().Error())
	default:
		if ok := s.inCache(cmd); ok {
			if err := s.cache.Remove(cmd); err != nil {
				s.logError(
					"cache",
					"Remove",
					err.Error(),
					fmt.Sprintf("remove BotCommand %q failed", cmd.GetCommand()),
				)
			}
		}

		if err := s.db.Remove(cmd); err != nil {
			s.logError(
				"db",
				"Remove",
				err.Error(),
				fmt.Sprintf("remove BotCommand %q failed", cmd.GetCommand()),
			)
			return &empty.Empty{}, status.Error(codes.Internal, "error while removing command")
		}
	}

	s.logInfo(
		"server",
		"DeleteCommand",
		fmt.Sprintf("BotCommand %q removed successfully", cmd.GetCommand()),
		time.Since(start),
	)
	return &empty.Empty{}, nil
}

func (s *server) inCache(cmd *proto.Command) bool {
	if _, err := s.cache.Get(cmd); err != nil {
		return false
	}

	return true
}
