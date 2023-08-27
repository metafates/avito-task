package server

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

// pg Returns a postgres connection
func (a *Server) pg() *pgx.Conn {
	return a.options.Connections.Postgres
}

// psql Returns a new SQL statement builder
func (a *Server) psql() sq.StatementBuilderType {
	return sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar)
}
