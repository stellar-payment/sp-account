package main

import (
	"github.com/stellar-payment/sp-account/cmd/webservice"
	"github.com/stellar-payment/sp-account/internal/component"
	"github.com/stellar-payment/sp-account/internal/config"
	"github.com/stellar-payment/sp-account/pkg/initutil"
)

var (
	buildVer  string = "unknown"
	buildTime string = "unknown"
)

func main() {
	config.Init(buildTime, buildVer)
	conf := config.Get()

	initutil.InitDirectory()

	logger := component.NewLogger(component.NewLoggerParams{
		ServiceName: conf.ServiceName,
		PrettyPrint: true,
	})

	webservice.Start(conf, logger)
}
