package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/repository"
	"github.com/stellar-payment/sp-account/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	// ----- Register
	RegisterCustomer(ctx context.Context, payload *dto.RegisterCustomerPayload) (err error)
	RegisterMerchant(ctx context.Context, payload *dto.RegisterMerchantPayload) (err error)

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
	HandleDeleteUser(ctx context.Context, params *indto.User) (err error)
}

type service struct {
	conf       *serviceConfig
	redis      *redis.Client
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
	Redis      *redis.Client
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
		redis:      params.Redis,
	}
}
