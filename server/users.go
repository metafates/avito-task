package server

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/server/api"
)

func (a *Server) userExists(ctx context.Context, id api.UserID) (bool, error) {
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

func (a *Server) createUser(ctx context.Context, id api.UserID) error {
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

func (a *Server) assignedSegments(ctx context.Context, id api.UserID) ([]api.UserSegment, error) {
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

	var segments []api.UserSegment
	for rows.Next() {
		var segment api.UserSegment

		if err = rows.Scan(&segment.Slug, &segment.Expires); err != nil {
			return nil, err
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

func (a *Server) userHasSegment(ctx context.Context, user api.UserID, segment api.Slug) (bool, error) {
	segments, err := a.assignedSegments(ctx, user)
	if err != nil {
		return false, err
	}

	for _, s := range segments {
		if s.Slug == segment {
			return true, nil
		}
	}

	return false, nil
}

func (a *Server) assignSegment(
	ctx context.Context,
	user api.UserID,
	segment api.Slug,
	assignment api.SegmentAssignment,
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

func (a *Server) depriveSegment(
	ctx context.Context,
	user api.UserID,
	segment api.Slug,
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
