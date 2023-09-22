package service

import (
	"time"

	"github.com/stellar-payment/sp-account/pkg/dto"
)

func (s *service) Ping() (pingResponse dto.PublicPingResponse) {
	return dto.PublicPingResponse{
		Message:   "KyaaNakaWaZettaiDame",
		Timestamp: time.Now().Unix(),
	}
}
