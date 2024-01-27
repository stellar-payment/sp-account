package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	// ----- Users
	FindUsers(ctx context.Context, params *indto.UserParams) (res []*indto.User, err error)
	CountUsers(ctx context.Context, params *indto.UserParams) (res int64, err error)
	FindUser(ctx context.Context, params *indto.UserParams) (res *indto.User, err error)
	InsertUser(ctx context.Context, payload *model.User) (err error)
	UpdateUser(ctx context.Context, payload *model.User) (err error)
	DeleteUser(ctx context.Context, params *indto.UserParams) (err error)

	// ----- Session
	FindSessionByToken(ctx context.Context, usr *indto.UserAccess) (res *indto.User, err error)
	InsertSessionToken(ctx context.Context, usr *indto.User, token string) (err error)
	InvalidateSessionToken(ctx context.Context, usr *indto.User) (err error)
}

type repository struct {
	db    *sqlx.DB
	redis *redis.Client
	conf  *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	DB      *sqlx.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		conf:  &repositoryConfig{},
		db:    params.DB,
		redis: params.Redis,
	}
}

var pgSquirrel = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
