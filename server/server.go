package server

import (
	"context"

	"github.com/metafates/avito-task/server/api"
)

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	options Options
}

// PostUsersId implements StrictServerInterface.
func (a *Server) PostUsersId(ctx context.Context, request api.PostUsersIdRequestObject) (api.PostUsersIdResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if exists {
		return api.PostUsersId409JSONResponse{
			Error: "user exists",
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
			Error: "segment not found",
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
			Error: "user not found",
		}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.DeleteUsersIdSegmentsSlug404JSONResponse{
			Error: "segment not found",
		}, nil
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
			Error: "user not found",
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
			Error: "segment exists",
		}, nil
	}

	if err = a.createSegment(ctx, request.Slug, api.SegmentCreation(*request.Body)); err != nil {
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
			Error: "user not found",
		}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return api.PostUsersIdSegmentsSlug404JSONResponse{
			Error: "segment not found",
		}, nil
	}

	hasSegment, err := a.userHasSegment(ctx, request.Id, request.Slug)
	if err != nil {
		return nil, err
	}

	if hasSegment {
		return api.PostUsersIdSegmentsSlug409JSONResponse{
			Error: "segment is already assigned",
		}, nil
	}

	if err = a.assignSegment(
		ctx,
		request.Id,
		request.Slug,
		api.SegmentAssignment(*request.Body),
	); err != nil {
		return nil, err
	}

	return api.PostUsersIdSegmentsSlug200Response{}, nil
}
