package scopeutil

import (
	"context"

	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/util/ctxutil"
)

func ValidateScope(ctx context.Context, roles ...int64) (ok bool) {
	userMeta, ok := ctxutil.GetCtx[*indto.User](ctx, inconst.AUTH_CTX_KEY)
	if !ok {
		return false
	}

	if userMeta.RoleID == 0 {
		return true
	}

	for _, v := range roles {
		if userMeta.RoleID == v {
			return true
		}
	}

	return false
}
