package server

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (a *Server) pg() *pgx.Conn {
	return a.options.Connections.Postgres
}

func (a *Server) psql() sq.StatementBuilderType {
	return sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar)
}
