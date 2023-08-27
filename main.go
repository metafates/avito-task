package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/metafates/avito-task/config"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	log.Logger.Info().Msg("loading config")
	cfg, err := config.Load(".")
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("config")
	}

	log.Logger.Info().Msg("connecting to databases")
	dbConnections, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("database")
	}

	log.Logger.Info().Msg("starting the server")
	err = server.Run(ctx, net.JoinHostPort("0.0.0.0", cfg.Port), server.Options{
		Connections: dbConnections,
	})

	if err != nil {
		log.Logger.Fatal().Err(err).Msg("server")
	}
}
