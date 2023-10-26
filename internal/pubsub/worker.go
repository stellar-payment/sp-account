package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/service"
)

type EventPubSub struct {
	logger  zerolog.Logger
	redis   *redis.Client
	service service.Service
}

type NewEventPubSubParams struct {
	Logger  zerolog.Logger
	Redis   *redis.Client
	Service service.Service
}

func NewEventPubSub(params *NewEventPubSubParams) *EventPubSub {
	return &EventPubSub{
		logger:  params.Logger,
		redis:   params.Redis,
		service: params.Service,
	}
}

func (pb *EventPubSub) Listen() {
	ctx := context.Background()

	subscriber := pb.redis.Subscribe(ctx, "")

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		switch msg.Channel {
		case inconst.TOPIC_DELETE_USER:
			data := &indto.User{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleDeleteUser(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		}
	}
}
