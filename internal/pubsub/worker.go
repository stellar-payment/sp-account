package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/service"
)

type EventPubSub struct {
	logger       zerolog.Logger
	redis        *redis.Client
	service      service.Service
	secureRoutes []string
}

type NewEventPubSubParams struct {
	Logger       zerolog.Logger
	Redis        *redis.Client
	Service      service.Service
	SecureRoutes []string
}

func NewEventPubSub(params *NewEventPubSubParams) *EventPubSub {
	return &EventPubSub{
		logger:       params.Logger,
		redis:        params.Redis,
		service:      params.Service,
		secureRoutes: params.SecureRoutes,
	}
}

func (pb *EventPubSub) Listen() {
	ctx := context.Background()

	subscriber := pb.redis.Subscribe(ctx,
		inconst.TOPIC_REQUEST_SECURE_ROUTE,
		inconst.TOPIC_DELETE_USER,
	)

	data := fmt.Sprintf("%s,%s", "account", strings.Join(pb.secureRoutes, ","))
	pb.redis.Publish(context.Background(), inconst.TOPIC_BROADCAST_SECURE_ROUTE, data)

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
		case inconst.TOPIC_REQUEST_SECURE_ROUTE:
			data := fmt.Sprintf("%s,%s", "payment", strings.Join(pb.secureRoutes, ","))
			pb.redis.Publish(context.Background(), inconst.TOPIC_BROADCAST_SECURE_ROUTE, data)
		}
	}
}
