package server

import (
	"context"
	"strings"
	"time"

	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server/api"
)

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	options Options
}

// GetAudit implements api.StrictServerInterface.
func (s *Server) GetAudit(ctx context.Context, request api.GetAuditRequestObject) (api.GetAuditResponseObject, error) {
	var filter auditFilter
	if from := request.Params.From; from != nil {
		filter.From = &from.Time
	}

	if to := request.Params.To; to != nil {
		filter.To = &to.Time
	}

	if user := request.Params.User; user != nil {
		filter.User = user
	}

	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	audit, err := s.audit(ctx, conn, filter)
	if err != nil {
		log.Logger.Err(err).Send()
		return nil, err
	}

	csv, err := audit.CSV()
	if err != nil {
		log.Logger.Err(err).Send()
		return nil, err
	}

	return api.GetAudit200TextcsvResponse{
		Body:          strings.NewReader(csv),
		ContentLength: int64(len(csv)),
	}, nil
}

// PostUsersId implements StrictServerInterface.
func (s *Server) PostUsersId(ctx context.Context, request api.PostUsersIdRequestObject) (api.PostUsersIdResponseObject, error) {
	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.userExists(ctx, conn, request.Id)
	if err != nil {
		return nil, err
	}

	if exists {
		return api.PostUsersId409JSONResponse{
			Error: errUserExists.Error(),
		}, nil
	}

	if err := s.createUser(ctx, conn, request.Id); err != nil {
		return nil, err
	}

	return api.PostUsersId201Response{}, nil
}

// DeleteSegmentsSlug implements StrictServerInterface.
func (s *Server) DeleteSegmentsSlug(ctx context.Context, request api.DeleteSegmentsSlugRequestObject) (api.DeleteSegmentsSlugResponseObject, error) {
	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.segmentExists(ctx, conn, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	if err = s.deleteSegment(ctx, conn, request.Slug); err != nil {
		return nil, err
	}

	return api.DeleteSegmentsSlug200Response{}, nil
}

// DeleteUsersIdSegmentsSlug implements StrictServerInterface.
func (s *Server) DeleteUsersIdSegmentsSlug(ctx context.Context, request api.DeleteUsersIdSegmentsSlugRequestObject) (api.DeleteUsersIdSegmentsSlugResponseObject, error) {
	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.userExists(ctx, conn, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	exists, err = s.segmentExists(ctx, conn, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	hasSegment, err := s.userHasSegment(ctx, conn, request.Id, request.Slug)
	if err != nil {
		return nil, err
	}

	if !hasSegment {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotAssigned.Error(),
		}, err
	}

	if err = s.depriveSegment(ctx, conn, request.Id, request.Slug); err != nil {
		return nil, err
	}

	return api.DeleteUsersIdSegmentsSlug200Response{}, nil
}

// GetUsersIdSegments implements StrictServerInterface.
func (s *Server) GetUsersIdSegments(ctx context.Context, request api.GetUsersIdSegmentsRequestObject) (api.GetUsersIdSegmentsResponseObject, error) {
	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.userExists(ctx, conn, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.GetUsersIdSegments404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	segments, err := s.assignedSegments(ctx, conn, request.Id)
	if err != nil {
		return nil, err
	}

	return api.GetUsersIdSegments200JSONResponse(segments), nil
}

// PostSegmentsSlug implements StrictServerInterface.
func (s *Server) PostSegmentsSlug(ctx context.Context, request api.PostSegmentsSlugRequestObject) (api.PostSegmentsSlugResponseObject, error) {
	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.segmentExists(ctx, conn, request.Slug)
	if err != nil {
		return nil, err
	}

	if exists {
		return api.PostSegmentsSlug409JSONResponse{
			Error: errSegmentExists.Error(),
		}, nil
	}

	var segmentCreation api.SegmentCreation
	if body := request.Body; body != nil {
		segmentCreation.Outreach = body.Outreach
	}

	if err = s.createSegment(ctx, conn, request.Slug, segmentCreation); err != nil {
		return nil, err
	}

	return api.PostSegmentsSlug201Response{}, nil
}

// PostUsersIdSegmentsSlug implements StrictServerInterface.
func (s *Server) PostUsersIdSegmentsSlug(ctx context.Context, request api.PostUsersIdSegmentsSlugRequestObject) (api.PostUsersIdSegmentsSlugResponseObject, error) {
	if expires := request.Body.ExpiresAt; expires != nil && expires.Before(time.Now()) {
		return api.PostUsersIdSegmentsSlug400JSONResponse{
			Error: "segment is already expired",
		}, nil
	}

	conn, err := s.pgpool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	exists, err := s.userExists(ctx, conn, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.PostUsersIdSegmentsSlug404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	exists, err = s.segmentExists(ctx, conn, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.PostUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	hasSegment, err := s.userHasSegment(ctx, conn, request.Id, request.Slug)
	if err != nil {
		return nil, err
	}

	if hasSegment {
		return api.PostUsersIdSegmentsSlug409JSONResponse{
			Error: errSegmentAssignedAlready.Error(),
		}, nil
	}

	var segmentAssignment api.SegmentAssignment
	if body := request.Body; body != nil {
		segmentAssignment.ExpiresAt = body.ExpiresAt
	}

	if err = s.assignSegment(
		ctx,
		conn,
		request.Id,
		request.Slug,
		segmentAssignment,
	); err != nil {
		return nil, err
	}

	return api.PostUsersIdSegmentsSlug200Response{}, nil
}
