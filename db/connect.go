package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metafates/avito-task/config"
	"github.com/metafates/avito-task/log"
)

//go:embed init.sql
var postgresInitSQL string

// Pools contains DB connections pools
type Pools struct {
	Postgres *pgxpool.Pool
}

// Connect initializes databases, establishes connections and pings them
func Connect(ctx context.Context, config config.DBConfig) (pools Pools, err error) {
	log.Logger.Info().Str("db", "postgres").Msg("connecting")
	pools.Postgres, err = pgxpool.New(ctx, config.PostgresURI)
	if err != nil {
		return
	}

	err = pools.Postgres.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		log.Logger.Info().Str("db", "postgres").Msg("ping")
		if err = c.Ping(ctx); err != nil {
			return err
		}
		log.Logger.Info().Str("db", "postgres").Msg("pong")

		log.Logger.Info().Str("db", "postgres").Msg("initializing")
		_, err = c.Exec(ctx, postgresInitSQL)
		if err != nil {
			return err
		}
		log.Logger.Info().Str("db", "postgres").Msg("initialized")
		return nil
	})

	if err != nil {
		return
	}

	return
}
