package api

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (a *API) pg() *pgx.Conn {
	return a.options.Connections.Postgres
}

func (a *API) psql() sq.StatementBuilderType {
	return sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar)
}
