package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-account/internal/util/echttputil"
	"github.com/stellar-payment/sp-account/pkg/dto"
)

type PingHandler func() (pingResponse dto.PublicPingResponse)

func HandlePing(handler PingHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := handler()
		return echttputil.WriteSuccessResponse(c, resp)
	}
}
