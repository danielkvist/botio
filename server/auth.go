package server

import (
	"context"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (s *server) jwtAuth(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := errors.New("error while getting metadata from incoming context")
		s.logError(
			"server",
			"jwtAuth",
			err.Error(),
			"error while getting metadata for JWT auth",
		)

		return nil, err
	}

	if len(md.Get("token")) == 0 {
		err := errors.New("metadata from incoming context no contains a token")
		s.logError(
			"server",
			"jwtAuth",
			err.Error(),
			"error while extracting from the metadata the JWT token for auth",
		)

		return nil, err
	}

	token, err := jwt.Parse(md.Get("token")[0], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("while parsing a JWT token there was an error")
		}

		return []byte(s.key), nil
	})

	if err != nil {
		s.logError(
			"server",
			"jwtAuth.jwtParse",
			err.Error(),
			"error while parsing the received token",
		)

		return nil, err
	}

	if !token.Valid {
		s.logError(
			"server",
			"jwtAuth",
			errors.New("no valid token received").Error(),
			"error while parsing the received token",
		)
		return nil, err
	}

	return ctx, nil
}
