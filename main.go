package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/metafates/avito-task/config"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	log.Logger.Info().Msg("loading config")
	if err := godotenv.Load(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Logger.Warn().Msg(".env file is missing")
		} else {
			log.Logger.Fatal().Err(err).Msg("config")
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("config")
	}

	log.Logger.Info().Msg("connecting to databases")
	pools, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("database")
	}

	log.Logger.Info().Msg("starting the server")
	err = server.Run(ctx, net.JoinHostPort("0.0.0.0", cfg.Port), server.Options{
		Pools: pools,
	})

	if err != nil {
		log.Logger.Fatal().Err(err).Msg("server")
	}
}
