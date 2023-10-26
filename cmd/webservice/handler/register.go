package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-account/internal/util/echttputil"
	"github.com/stellar-payment/sp-account/pkg/dto"
	"github.com/stellar-payment/sp-account/pkg/errs"
)

type RegisterCustomerHandler func(context.Context, *dto.RegisterCustomerPayload) error

func HandleRegisterCustomer(handler RegisterCustomerHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.RegisterCustomerPayload{}
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

type RegisterMerchantHandler func(context.Context, *dto.RegisterMerchantPayload) error

func HandleRegisterMerchant(handler RegisterMerchantHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.RegisterMerchantPayload{}
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
