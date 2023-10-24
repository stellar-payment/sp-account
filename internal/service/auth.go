package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/model"
	"github.com/stellar-payment/sp-account/internal/util/sessionutil"
	"github.com/stellar-payment/sp-account/internal/util/structutil"
	"github.com/stellar-payment/sp-account/pkg/dto"
	"github.com/stellar-payment/sp-account/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) AuthSignup(ctx context.Context, payload *dto.UserPayload) (err error) {
	logger := zerolog.Ctx(ctx)

	if val := structutil.CheckMandatoryField(payload); val != "" {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	if usr, err := s.repository.FindUser(ctx, &indto.UserParams{Username: payload.Username}); err != nil {
		logger.Error().Err(err).Msg("failed to lookup existing user")
		return err
	} else if usr != nil {
		return errs.ErrBrokenUserReq
	}

	if payload.RoleID != inconst.ROLE_CUSTOMER && payload.RoleID != inconst.ROLE_MERCHANT {
		payload.RoleID = inconst.ROLE_CUSTOMER
	}

	saltedPass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return err
	}

	user := &model.User{
		UserID:   uuid.NewString(),
		Username: payload.Username,
		Password: string(saltedPass),
		RoleID:   payload.RoleID,
	}

	err = s.repository.InsertUser(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to insert new user")
		return
	}

	return
}

func (s *service) AuthLogin(ctx context.Context, payload *dto.AuthLoginPayload) (res *dto.AuthResponse, err error) {
	logger := log.Ctx(ctx)

	if val := structutil.CheckMandatoryField(payload); val != "" {
		logger.Error().Msgf("field %s is missing a value", val)
		return nil, errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	user, err := s.repository.FindUser(ctx, &indto.UserParams{Username: payload.Username})
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch userdata")
		return
	} else if user == nil {
		return nil, errs.ErrNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		logger.Error().Err(err).Msg("failed to validate user credentials")
		return nil, errs.ErrInvalidCred
	}

	token := sessionutil.GenerateSessionKey()
	err = s.repository.InvalidateSessionToken(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to invalidate user session")
		return nil, err
	}

	err = s.repository.InsertSessionToken(ctx, user, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed to insert user session")
		return nil, err
	}

	res = &dto.AuthResponse{
		UserID:      user.UserID,
		Username:    user.Username,
		RoleID:      user.RoleID,
		AccessToken: token,
	}

	return
}

// func (s *service) newTempToken(ctx context.Context, user *indto.User) (signed string, err error) {
// 	conf := config.Get()
// 	logger := log.Ctx(ctx)

// 	claims := s.newATUserClaim(user.UserID, user.Username, 5*time.Minute)
// 	accessToken := jwt.NewWithClaims(conf.JWT_SIGNING_METHOD, claims)
// 	signed, err = accessToken.SignedString(conf.JWT_SIGNATURE_KEY)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("failed to signed new access token")
// 		return "", err
// 	}

// 	return
// }
