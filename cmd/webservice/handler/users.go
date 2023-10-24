package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-account/internal/util/echttputil"
	"github.com/stellar-payment/sp-account/pkg/dto"
	"github.com/stellar-payment/sp-account/pkg/errs"
)

type GetUsersHandler func(context.Context, *dto.UsersQueryParams) (*dto.ListUserResponse, error)

func HandleGetUsers(handler GetUsersHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.UsersQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		res, err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetUserByIDHandler func(context.Context, *dto.UsersQueryParams) (*dto.UserResponse, error)

func HandleGetUserByID(handler GetUserByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.UsersQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		res, err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetUserMeHandler func(context.Context) (*dto.UserResponse, error)

func HandleGetUserMe(handler GetUserMeHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type CreateUserHandler func(context.Context, *dto.UserPayload) error

func HandleCreateUsers(handler CreateUserHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.UserPayload{}
		if err := c.Bind(payload); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), payload)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type UpdateUserHandler func(context.Context, *dto.UsersQueryParams, *dto.UserPayload) error

func HandleUpdateUsers(handler UpdateUserHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.UsersQueryParams{
			UserID: c.Param("userID"),
		}

		payload := &dto.UserPayload{}
		if err := c.Bind(payload); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), params, payload)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type DeleteUserHandler func(context.Context, *dto.UsersQueryParams) error

func HandleDeleteUser(handler DeleteUserHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.UsersQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
