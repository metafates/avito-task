package server

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/server/api"
)

// segmentWithOutreach Segment that is guaranteed to have non-nil Outreach column
type segmentWithOutreach struct {
	Slug     api.Slug
	Outreach api.Outreach
}

// shouldSegmentBeAssigned Determines if the segment shold be assigned
// to some user based on the segments outreach.
//
// This function is not determenistic. That is, it may return different
// result each time
func shouldSegmentBeAssigned(outreach api.Outreach) bool {
	return rand.Float32() <= outreach
}

// segmentExists Checks if there is a segment with the given slug
func (s *Server) segmentExists(ctx context.Context, conn *pgxpool.Conn, slug api.Slug) (bool, error) {
	sql, args, err := s.
		psql().
		Select("1").
		From(db.TableSegments).
		Where(squirrel.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return false, err
	}

	// postgres specifc
	row := conn.QueryRow(ctx, fmt.Sprintf("SELECT EXISTS(%s)", sql), args...)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// createSegment Creates a new segment.
func (s *Server) createSegment(ctx context.Context, conn *pgxpool.Conn, slug api.Slug, segment api.SegmentCreation) error {
	colums := []string{"slug"}
	values := []any{slug}

	if segment.Outreach != nil {
		colums = append(colums, "outreach")
		values = append(values, segment.Outreach)
	}

	sql, args, err := s.
		psql().
		Insert(db.TableSegments).
		Columns(colums...).
		Values(values...).
		ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if segment.Outreach == nil {
		return nil
	}
	outreach := *segment.Outreach

	// assign this segment to the users based on its outreach
	users, err := s.users(ctx, conn)
	if err != nil {
		return err
	}

	for _, user := range users {
		if !shouldSegmentBeAssigned(outreach) {
			continue
		}

		if err = s.assignSegment(ctx, conn, user, slug, api.SegmentAssignment{}); err != nil {
			return err
		}
	}

	return nil
}

// deleteSegment Deletes the segment
func (s *Server) deleteSegment(ctx context.Context, conn *pgxpool.Conn, slug api.Slug) error {
	sql, args, err := s.
		psql().
		Delete(db.TableSegments).
		Where(squirrel.Eq{"slug": slug}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

// segmentsWithOutreach Returns a list of segments that has non-nil outreach column
func (s *Server) segmentsWithOutreach(ctx context.Context, conn *pgxpool.Conn) ([]segmentWithOutreach, error) {
	sql, args, err := s.
		psql().
		Select("slug", "outreach").
		From(db.TableSegments).
		Where(squirrel.NotEq{"outreach": nil}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var segments []segmentWithOutreach
	for rows.Next() {
		var segment segmentWithOutreach

		if err = rows.Scan(&segment.Slug, &segment.Outreach); err != nil {
			return nil, err
		}

		segments = append(segments, segment)
	}

	return segments, nil
}
