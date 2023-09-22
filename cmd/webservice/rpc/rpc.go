package rpc

import (
	hentaiRPC "github.com/nmluci/gostellar/pkg/rpc/hentai"
	insvc "github.com/stellar-payment/sp-account/internal/service"
)

type HentaiRPC struct {
	hentaiRPC.UnimplementedNakaZettaiDameServer
	service insvc.Service
}

func Init(svc insvc.Service) hentaiRPC.NakaZettaiDameServer {
	return &HentaiRPC{
		service: svc,
	}
}
