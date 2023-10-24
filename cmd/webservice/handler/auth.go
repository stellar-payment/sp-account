package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-account/internal/util/echttputil"
	"github.com/stellar-payment/sp-account/pkg/dto"
)

type SignupHandler func(ctx context.Context, req *dto.UserPayload) (err error)

func HandleSignup(handler SignupHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.UserPayload{}
		if err := c.Bind(req); err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type AuthLoginHandler func(ctx context.Context, req *dto.AuthLoginPayload) (res *dto.AuthResponse, err error)

func HandleAuthLogin(handler AuthLoginHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.AuthLoginPayload{}
		if err := c.Bind(req); err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		res, err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}
