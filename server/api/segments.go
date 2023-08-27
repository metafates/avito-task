package api

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/metafates/avito-task/db"
)

func (a *API) segmentExists(ctx context.Context, slug Slug) (bool, error) {
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

func (a *API) createSegment(ctx context.Context, slug Slug, segment SegmentCreation) error {
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

func (a *API) deleteSegment(ctx context.Context, slug Slug) error {
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
