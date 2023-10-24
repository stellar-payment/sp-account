package pubsub

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/service"
)

var (
	tagLoggerPBListen = "[PubSub-Listen]"
)

type FilePubSub struct {
	logger  zerolog.Logger
	redis   *redis.Client
	service service.Service
}

type NewFilePubSubParams struct {
	Logger  zerolog.Logger
	Redis   *redis.Client
	Service service.Service
}

func NewFileSub(params NewFilePubSubParams) *FilePubSub {
	return &FilePubSub{
		logger:  params.Logger,
		redis:   params.Redis,
		service: params.Service,
	}
}

func (pb *FilePubSub) Listen() {
	ctx := context.Background()

	subscriber := pb.redis.Subscribe(ctx, "")

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		fmt.Print(msg)
	}
}
