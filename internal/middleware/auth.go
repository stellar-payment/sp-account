package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/service"
	"github.com/stellar-payment/sp-account/internal/util/ctxutil"
	"github.com/stellar-payment/sp-account/internal/util/echttputil"
	"github.com/stellar-payment/sp-account/pkg/errs"
)

func AuthorizationMiddleware(svc service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header

			var err error
			token := header.Get("Authorization")

			if token == "" {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			splittedToken := strings.Split(token, " ")
			if len(splittedToken) != 2 {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			if splittedToken[0] != "Bearer" {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			accessToken := splittedToken[1]
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			ctx, err := svc.AuthorizedAccessCtx(c.Request().Context(), accessToken)
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	}
}

func OptionalAuthorizationMiddleware(svc service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header

			var err error
			token := header.Get("Authorization")

			if token == "" {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			splittedToken := strings.Split(token, " ")
			if len(splittedToken) != 2 {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			if splittedToken[0] != "Bearer" {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			accessToken := splittedToken[1]
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			ctx := ctxutil.WrapCtx(c.Request().Context(), inconst.AT_CTX_KEY, accessToken)
			authCtx, _ := svc.AuthorizedAccessCtx(ctx, accessToken)

			c.SetRequest(c.Request().Clone(authCtx))
			return next(c)
		}
	}
}
