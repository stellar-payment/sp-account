package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/config"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/pkg/errs"
)

var (
	SESSION_TOKEN = "auth:session:%s"
	USER_TOKEN    = "auth:user:%s:session"
)

func (r *repository) FindSessionByToken(ctx context.Context, usr *indto.UserAccess) (res *indto.User, err error) {
	logger := zerolog.Ctx(ctx)

	data, err := r.redis.Get(ctx, fmt.Sprintf(SESSION_TOKEN, usr.AccessToken)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errs.ErrInvalidCred
	} else if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return nil, err
	}

	res = &indto.User{UserID: data}
	if err != nil {
		logger.Error().Msgf("parser err: %+v", err)
		return nil, err
	}

	return
}

func (r *repository) InsertSessionToken(ctx context.Context, usr *indto.User, token string) (err error) {
	conf := config.Get()
	logger := zerolog.Ctx(ctx)

	err = r.redis.Set(ctx, fmt.Sprintf(SESSION_TOKEN, token), usr.UserID, conf.RedisConfig.DefaultExp).Err()
	if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return err
	}

	err = r.redis.Set(ctx, fmt.Sprintf(USER_TOKEN, usr.UserID), token, conf.RedisConfig.DefaultExp).Err()
	if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return err
	}

	return
}

func (r *repository) InvalidateSessionToken(ctx context.Context, usr *indto.User) (err error) {
	logger := zerolog.Ctx(ctx)

	activeToken, err := r.redis.Get(ctx, fmt.Sprintf(USER_TOKEN, usr.UserID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return err
	}

	err = r.redis.Del(ctx, fmt.Sprintf(SESSION_TOKEN, activeToken)).Err()
	if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return err
	}

	err = r.redis.Del(ctx, fmt.Sprintf(USER_TOKEN, usr.UserID)).Err()
	if err != nil {
		logger.Error().Msgf("redis err: %+v", err)
		return err
	}

	return
}
