package server

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pg Returns a postgres connection pool
func (s *Server) pgpool() *pgxpool.Pool {
	return s.options.Pools.Postgres
}

// psql Returns a new SQL statement builder
func (s *Server) psql() squirrel.StatementBuilderType {
	return squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar)
}
