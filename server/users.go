package server

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server/api"
)

// users Returns a list of all users
func (a *Server) users(ctx context.Context) ([]api.UserID, error) {
	sql, args, err := a.psql().Select("id").From(db.TableUsers).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.pg().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var users []api.UserID
	for rows.Next() {
		var user api.UserID

		if err = rows.Scan(&user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// userExists Checks if the given user exists
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

// createUser Creates a new user.
func (a *Server) createUser(ctx context.Context, id api.UserID) error {
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
	if err != nil {
		return err
	}

	// Randomly assign segments based on their outreach
	segments, err := a.segmentsWithOutreach(ctx)
	if err != nil {
		return err
	}

	for _, segment := range segments {
		if !shouldSegmentBeAssigned(segment.Outreach) {
			continue
		}

		if err = a.assignSegment(ctx, id, segment.Slug, api.SegmentAssignment{}); err != nil {
			return err
		}
	}

	return nil
}

// assignedSegments Gets all segments assigned to the user
func (a *Server) assignedSegments(ctx context.Context, id api.UserID) ([]api.UserSegment, error) {
	sql, args, err := a.
		psql().
		Select("segment_slug", "expires_at").
		From(db.TableAssignedSegments).
		Where(squirrel.Eq{"user_id": id}).
		Where(squirrel.Or{
			squirrel.Gt{"expires_at": time.Now()},
			squirrel.Eq{"expires_at": nil},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(sql)

	rows, err := a.pg().Query(ctx, sql, args...)
	if err != nil {
		log.Logger.Err(err).Send()
		return nil, err
	}

	segments := make([]api.UserSegment, 0)
	for rows.Next() {
		var segment api.UserSegment

		if err = rows.Scan(&segment.Slug, &segment.ExpiresAt); err != nil {
			return nil, err
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// userHasSegment Checks if user is assigned to this segment
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

// assignSegment Assigns segment to a user
func (a *Server) assignSegment(
	ctx context.Context,
	user api.UserID,
	segment api.Slug,
	assignment api.SegmentAssignment,
) error {
	colums := []string{"user_id", "segment_slug"}
	values := []any{user, segment}

	if assignment.ExpiresAt != nil {
		colums = append(colums, "expires_at")
		values = append(values, assignment.ExpiresAt)
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

// depriveSegment Removes segment from a user
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
