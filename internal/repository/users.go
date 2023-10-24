package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/model"
)

func (r *repository) FindUsers(ctx context.Context, params *indto.UserParams) (res []*indto.User, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"u.deleted_at": nil},
	}

	baseStmt := pgSquirrel.Select("u.id", "u.username", "u.password", "u.role_id").From("users u").Where(cond)

	if params.Limit != 0 && params.Page >= 1 {
		baseStmt = baseStmt.Limit(params.Limit).Offset((params.Page - 1) * params.Limit)
	}

	stmt, args, err := baseStmt.ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = []*indto.User{}
	rows, err := r.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	for rows.Next() {
		temp := &indto.User{}

		if err = rows.StructScan(temp); err != nil {
			logger.Error().Err(err).Msg("sql map err")
			return
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) CountUsers(ctx context.Context, params *indto.UserParams) (res int64, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"u.deleted_at": nil},
	}

	stmt, args, err := pgSquirrel.Select("count(*)").From("users u").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	err = r.db.QueryRowxContext(ctx, stmt, args...).Scan(&res)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}

func (r *repository) FindUser(ctx context.Context, params *indto.UserParams) (res *indto.User, err error) {
	logger := zerolog.Ctx(ctx)

	cond := squirrel.And{
		squirrel.Eq{"u.deleted_at": nil},
	}

	if params.Username != "" {
		cond = append(cond, squirrel.Eq{"username": params.Username})
	} else if params.UserID != "" {
		cond = append(cond, squirrel.Eq{"id": params.UserID})
	}

	stmt, args, err := pgSquirrel.Select("u.id", "u.username", "u.password", "u.role_id").From("users u").Where(cond).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	res = &indto.User{}
	err = r.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("sql err")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (r *repository) InsertUser(ctx context.Context, payload *model.User) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Insert("users").Columns("id", "username", "password", "role_id").
		Values(payload.UserID, payload.Username, payload.Password, payload.RoleID).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}

func (r *repository) UpdateUser(ctx context.Context, payload *model.User) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("users").SetMap(map[string]interface{}{
		"username":   squirrel.Expr("coalesce(nullif(?, ''), username)", payload.Username),
		"password":   squirrel.Expr("coalesce(nullif(?, ''), password)", payload.Password),
		"updated_at": time.Now(),
	}).Where(squirrel.And{
		squirrel.Eq{"id": payload.UserID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}
	fmt.Println(stmt)
	fmt.Println(args)

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}

func (r *repository) DeleteUser(ctx context.Context, params *indto.UserParams) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, err := pgSquirrel.Update("users").SetMap(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
	}).Where(squirrel.And{
		squirrel.Eq{"id": params.UserID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()
	if err != nil {
		logger.Error().Err(err).Msg("squirrel err")
		return
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("sql err")
		return
	}

	return
}
