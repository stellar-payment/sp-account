package service

import (
	"context"

	"github.com/stellar-payment/sp-account/internal/repository"
	"github.com/stellar-payment/sp-account/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	// ----- Auth
	AuthSignup(ctx context.Context, payload *dto.UserPayload) (err error)
	AuthLogin(ctx context.Context, payload *dto.AuthLoginPayload) (res *dto.AuthResponse, err error)
	AuthorizedAccessCtx(ctx context.Context, token string) (res context.Context, err error)

	// ----- Users
	GetAllUser(ctx context.Context, params *dto.UsersQueryParams) (res *dto.ListUserResponse, err error)
	GetUser(ctx context.Context, params *dto.UsersQueryParams) (res *dto.UserResponse, err error)
	GetUserMe(ctx context.Context) (res *dto.UserResponse, err error)
	CreateUser(ctx context.Context, payload *dto.UserPayload) (err error)
	UpdateUser(ctx context.Context, params *dto.UsersQueryParams, payload *dto.UserPayload) (err error)
	DeleteUser(ctx context.Context, params *dto.UsersQueryParams) (err error)
}

type service struct {
	conf       *serviceConfig
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
	}
}
