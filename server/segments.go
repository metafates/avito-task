package server

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/server/api"
)

func (a *Server) segmentExists(ctx context.Context, slug api.Slug) (bool, error) {
	sql, args, err := a.
		psql().
		Select("1").
		From(db.TableSegments).
		Where(squirrel.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return false, err
	}

	// postgres specifc
	row := a.pg().QueryRow(ctx, fmt.Sprintf("SELECT EXISTS(%s)", sql), args...)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (a *Server) createSegment(ctx context.Context, slug api.Slug, segment api.SegmentCreation) error {
	colums := []string{"slug"}
	values := []any{slug}

	if segment.Outreach != nil {
		colums = append(colums, "outreach")
		values = append(values, segment.Outreach)
	}

	sql, args, err := a.
		psql().
		Insert(db.TableSegments).
		Columns(colums...).
		Values(values...).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}

func (a *Server) deleteSegment(ctx context.Context, slug api.Slug) error {
	sql, args, err := a.
		psql().
		Delete(db.TableSegments).
		Where(squirrel.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}
