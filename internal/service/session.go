package service

import (
	"context"
)

var (
	tagLoggerAuthenticateSession = "[AuthenticateSession]"
	tagLoggerAuthenticateService = "[AuthenticateService]"
)

func (s *service) AuthenticateSession(ctx context.Context, token string) (access context.Context, err error) {
	return
}

func (s *service) AuthenticateService(ctx context.Context, name string) (access context.Context, err error) {
	return
}
