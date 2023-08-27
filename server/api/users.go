package api

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/metafates/avito-task/db"
)

func (a *API) userExists(ctx context.Context, id UserID) (bool, error) {
	sql, args, err := a.
		psql().
		Select("1").
		From(db.TableUsers).
		Where(squirrel.Eq{"id": id}).
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

func (a *API) createUser(ctx context.Context, id UserID) error {
	// automatically assign segments to a user based on their outreach
	sql, args, err := a.
		psql().
		Insert(db.TableUsers).
		Columns("id").
		Values(id).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}

func (a *API) assignedSegments(ctx context.Context, id UserID) ([]UserSegment, error) {
	sql, args, err := a.
		psql().
		Select("segment_slug", "expires_at").
		From(db.TableAssignedSegments).
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.pg().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var segments []UserSegment
	for rows.Next() {
		var segment UserSegment

		if err = rows.Scan(&segment.Slug, &segment.Expires); err != nil {
			return nil, err
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

func (a *API) assignSegment(
	ctx context.Context,
	user UserID,
	segment Slug,
	assignment SegmentAssignment,
) error {
	colums := []string{"user_id", "segment_slug"}
	values := []any{user, segment}

	if assignment.Expires != nil {
		colums = append(colums, "expires_at")
		values = append(values, assignment.Expires)
	}

	sql, args, err := a.
		psql().
		Insert(db.TableAssignedSegments).
		Columns(colums...).
		Values(values...).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}

func (a *API) depriveSegment(
	ctx context.Context,
	user UserID,
	segment Slug,
) error {
	sql, args, err := a.
		psql().
		Delete(db.TableAssignedSegments).
		Where(squirrel.Eq{"user_id": user, "segment_slug": segment}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}
