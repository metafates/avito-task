package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/metafates/avito-task/config"
	"github.com/metafates/avito-task/log"
)

//go:embed init.sql
var postgresInitSQL string

// Connections is a tuple of DB connections
type Connections struct {
	Postgres *pgx.Conn
}

// Connect initializes databases, establishes connections and pings them
func Connect(ctx context.Context, config config.DBConfig) (connections Connections, err error) {
	log.Logger.Info().Str("db", "postgres").Msg("connecting")
	connections.Postgres, err = pgx.Connect(ctx, config.PostgresURI)
	if err != nil {
		return
	}

	log.Logger.Info().Str("db", "postgres").Msg("ping")
	if err = connections.Postgres.Ping(ctx); err != nil {
		return
	}
	log.Logger.Info().Str("db", "postgres").Msg("pong")

	log.Logger.Info().Str("db", "postgres").Msg("initializing")
	_, err = connections.Postgres.Exec(ctx, postgresInitSQL)
	if err != nil {
		return
	}
	log.Logger.Info().Str("db", "postgres").Msg("initialized")

	return
}
