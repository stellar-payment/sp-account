package webservice

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/cmd/webservice/router"
	"github.com/stellar-payment/sp-account/internal/component"
	"github.com/stellar-payment/sp-account/internal/config"
	"github.com/stellar-payment/sp-account/internal/repository"
	"github.com/stellar-payment/sp-account/internal/service"
)

func Start(conf *config.Config, logger zerolog.Logger) {
	db, err := component.InitPostgres(&component.InitPostgresParams{
		Conf:   &conf.PostgresConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize db")
	}

	redis, err := component.InitRedis(&component.InitRedisParams{
		Conf:   &conf.RedisConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize redis")
	}

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	repo := repository.NewRepository(&repository.NewRepositoryParams{
		DB:    db,
		Redis: redis,
	})

	service := service.NewService(&service.NewServiceParams{
		Repository: repo,
	})

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	logger.Info().Msgf("starting service, listening to: %s", conf.ServiceAddress)

	if err := ec.Start(conf.ServiceAddress); err != nil {
		logger.Error().Msgf("starting service, cause: %+v", err)
	}
}
