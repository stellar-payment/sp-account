package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/util/ctxutil"
	"github.com/stellar-payment/sp-account/pkg/errs"
)

func (s *service) AuthorizedAccessCtx(ctx context.Context, token string) (res context.Context, err error) {
	logger := zerolog.Ctx(ctx)

	usr, err := s.repository.FindSessionByToken(ctx, &indto.UserAccess{AccessToken: token})
	if err != nil {
		logger.Error().Err(err).Msg("failed to authorize access token")
		return ctx, err
	}

	usrdata, err := s.repository.FindUser(ctx, &indto.UserParams{UserID: usr.UserID})
	if err != nil {
		logger.Error().Err(err).Msg("failed to lookup user")
		return ctx, err
	} else if usr == nil {
		logger.Error().Err(err).Msg("user not found")
		return ctx, errs.ErrInvalidCred
	}

	res = ctxutil.WrapCtx(ctx, inconst.AUTH_CTX_KEY, usrdata)
	return
}
