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
	var from, to *time.Time
	if request.Params.From != nil {
		from = &request.Params.From.Time
	}

	if request.Params.To != nil {
		to = &request.Params.To.Time
	}

	audit, err := s.audit(ctx, from, to)
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
func (a *Server) PostUsersId(ctx context.Context, request api.PostUsersIdRequestObject) (api.PostUsersIdResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if exists {
		return api.PostUsersId409JSONResponse{
			Error: errUserExists.Error(),
		}, nil
	}

	if err := a.createUser(ctx, request.Id); err != nil {
		return nil, err
	}

	return api.PostUsersId201Response{}, nil
}

// DeleteSegmentsSlug implements StrictServerInterface.
func (a *Server) DeleteSegmentsSlug(ctx context.Context, request api.DeleteSegmentsSlugRequestObject) (api.DeleteSegmentsSlugResponseObject, error) {
	exists, err := a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	if err = a.deleteSegment(ctx, request.Slug); err != nil {
		return nil, err
	}

	return api.DeleteSegmentsSlug200Response{}, nil
}

// DeleteUsersIdSegmentsSlug implements StrictServerInterface.
func (a *Server) DeleteUsersIdSegmentsSlug(ctx context.Context, request api.DeleteUsersIdSegmentsSlugRequestObject) (api.DeleteUsersIdSegmentsSlugResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	hasSegment, err := a.userHasSegment(ctx, request.Id, request.Slug)
	if err != nil {
		return nil, err
	}

	if !hasSegment {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotAssigned.Error(),
		}, err
	}

	if err = a.depriveSegment(ctx, request.Id, request.Slug); err != nil {
		return nil, err
	}

	return api.DeleteUsersIdSegmentsSlug200Response{}, nil
}

// GetUsersIdSegments implements StrictServerInterface.
func (a *Server) GetUsersIdSegments(ctx context.Context, request api.GetUsersIdSegmentsRequestObject) (api.GetUsersIdSegmentsResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.GetUsersIdSegments404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	segments, err := a.assignedSegments(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return api.GetUsersIdSegments200JSONResponse(segments), nil
}

// PostSegmentsSlug implements StrictServerInterface.
func (a *Server) PostSegmentsSlug(ctx context.Context, request api.PostSegmentsSlugRequestObject) (api.PostSegmentsSlugResponseObject, error) {
	exists, err := a.segmentExists(ctx, request.Slug)
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

	if err = a.createSegment(ctx, request.Slug, segmentCreation); err != nil {
		return nil, err
	}

	return api.PostSegmentsSlug201Response{}, nil
}

// PostUsersIdSegmentsSlug implements StrictServerInterface.
func (a *Server) PostUsersIdSegmentsSlug(ctx context.Context, request api.PostUsersIdSegmentsSlugRequestObject) (api.PostUsersIdSegmentsSlugResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.PostUsersIdSegmentsSlug404JSONResponse{
			Error: errUserNotFound.Error(),
		}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.PostUsersIdSegmentsSlug404JSONResponse{
			Error: errSegmentNotFound.Error(),
		}, nil
	}

	hasSegment, err := a.userHasSegment(ctx, request.Id, request.Slug)
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
		segmentAssignment.Expires = body.Expires
	}

	if err = a.assignSegment(
		ctx,
		request.Id,
		request.Slug,
		segmentAssignment,
	); err != nil {
		return nil, err
	}

	return api.PostUsersIdSegmentsSlug200Response{}, nil
}
