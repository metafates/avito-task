package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/metafates/avito-task/config"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server"
	"github.com/metafates/avito-task/server/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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
	err = server.Run(ctx, net.JoinHostPort("0.0.0.0", cfg.Port), api.Options{
		RefreshTokenDuration: 7 * 24 * time.Hour,
		AccessTokenDuration:  5 * time.Minute,
		Secret:               []byte(cfg.JWT.Secret),
		SigningMethod:        jwt.SigningMethodHS512,
		DB:                   dbConnections,
	})

	if err != nil {
		log.Logger.Fatal().Err(err).Msg("server")
	}
}
