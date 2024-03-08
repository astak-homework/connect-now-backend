package main

import (
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/astak-homework/connect-now-backend/server"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	app := server.NewApp(cfg)

	if err := app.Run(cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
