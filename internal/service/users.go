package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/model"
	"github.com/stellar-payment/sp-account/internal/util/ctxutil"
	"github.com/stellar-payment/sp-account/internal/util/scopeutil"
	"github.com/stellar-payment/sp-account/internal/util/structutil"
	"github.com/stellar-payment/sp-account/pkg/dto"
	"github.com/stellar-payment/sp-account/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) GetAllUser(ctx context.Context, params *dto.UsersQueryParams) (res *dto.ListUserResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 100
	}

	repoParams := &indto.UserParams{
		Keyword: params.Keyword,
		Limit:   params.Limit,
		Page:    params.Page,
	}

	res = &dto.ListUserResponse{
		Users: []*dto.UserResponse{},
		Meta: dto.ListPaginations{
			Limit: params.Limit,
			Page:  params.Page,
		},
	}

	count, err := s.repository.CountUsers(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return
	}

	res.Meta.TotalItem = uint64(count)
	res.Meta.TotalPage = uint64(math.Ceil(float64(count) / float64(params.Limit)))

	data, err := s.repository.FindUsers(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, v := range data {
		temp := &dto.UserResponse{
			UserID:   v.UserID,
			Username: v.Username,
			RoleID:   0,
		}

		res.Users = append(res.Users, temp)
	}

	return
}

func (s *service) GetUser(ctx context.Context, params *dto.UsersQueryParams) (res *dto.UserResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return nil, errs.ErrNoAccess
	}

	data, err := s.repository.FindUser(ctx, &indto.UserParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.UserResponse{
		UserID:   data.UserID,
		Username: data.Username,
		RoleID:   data.RoleID,
	}

	return
}

func (s *service) GetUserMe(ctx context.Context) (res *dto.UserResponse, err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN, inconst.ROLE_CUSTOMER, inconst.ROLE_MERCHANT); !ok {
		return nil, errs.ErrNoAccess
	}

	ctxmeta := ctxutil.GetUserCTX(ctx)

	data, err := s.repository.FindUser(ctx, &indto.UserParams{UserID: ctxmeta.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if data == nil {
		return nil, errs.ErrNotFound
	}

	res = &dto.UserResponse{
		UserID:   data.UserID,
		Username: data.Username,
		RoleID:   data.RoleID,
	}

	return
}

func (s *service) CreateUser(ctx context.Context, payload *dto.UserPayload) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	if val := structutil.CheckMandatoryField(payload); val != "" {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	if exists, err := s.repository.FindUser(ctx, &indto.UserParams{Username: payload.Username}); err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if exists != nil {
		return errs.ErrDuplicatedResources
	}

	userModel := &model.User{
		UserID:   uuid.NewString(),
		Username: payload.Username,
		Password: "",
		RoleID:   0,
	}

	if hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost); err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return err
	} else {
		userModel.Password = string(hashed)
	}

	err = s.repository.InsertUser(ctx, userModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) UpdateUser(ctx context.Context, params *dto.UsersQueryParams, payload *dto.UserPayload) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	if payload.Username != "" {
		if exists, err := s.repository.FindUser(ctx, &indto.UserParams{Username: payload.Username}); err != nil {
			logger.Error().Err(err).Send()
			return err
		} else if exists != nil && exists.UserID != params.UserID {
			return errs.ErrDuplicatedResources
		}
	}

	if payload.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error().Err(err).Msg("failed to hash password")
			return err
		}

		payload.Password = string(hashed)
	}

	userModel := &model.User{
		UserID:   params.UserID,
		Username: payload.Username,
		Password: payload.Password,
		RoleID:   payload.RoleID,
	}
	if err = s.repository.UpdateUser(ctx, userModel); err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) DeleteUser(ctx context.Context, params *dto.UsersQueryParams) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteUser(ctx, &indto.UserParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (s *service) HandleDeleteUser(ctx context.Context, params *indto.User) (err error) {
	logger := log.Ctx(ctx)

	if ok := scopeutil.ValidateScope(ctx, inconst.ROLE_ADMIN); !ok {
		return errs.ErrNoAccess
	}

	err = s.repository.DeleteUser(ctx, &indto.UserParams{UserID: params.UserID})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
