package main

import (
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/astak-homework/connect-now-backend/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
