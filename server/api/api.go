package api

import (
	"context"
)

var _ StrictServerInterface = (*API)(nil)

type API struct {
	options Options
}

// PostUsersId implements StrictServerInterface.
func (a *API) PostUsersId(ctx context.Context, request PostUsersIdRequestObject) (PostUsersIdResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if exists {
		return PostUsersId409Response{}, nil
	}

	if err := a.createUser(ctx, request.Id); err != nil {
		return nil, err
	}

	return PostUsersId201Response{}, nil
}

// DeleteSegmentsSlug implements StrictServerInterface.
func (a *API) DeleteSegmentsSlug(ctx context.Context, request DeleteSegmentsSlugRequestObject) (DeleteSegmentsSlugResponseObject, error) {
	exists, err := a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return DeleteSegmentsSlug404Response{}, nil
	}

	if err = a.deleteSegment(ctx, request.Slug); err != nil {
		return nil, err
	}

	return DeleteSegmentsSlug200Response{}, nil
}

// DeleteUsersIdSegmentsSlug implements StrictServerInterface.
func (a *API) DeleteUsersIdSegmentsSlug(ctx context.Context, request DeleteUsersIdSegmentsSlugRequestObject) (DeleteUsersIdSegmentsSlugResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return DeleteUsersIdSegmentsSlug404Response{}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return DeleteUsersIdSegmentsSlug404Response{}, nil
	}

	if err = a.depriveSegment(ctx, request.Id, request.Slug); err != nil {
		return nil, err
	}

	return DeleteUsersIdSegmentsSlug200Response{}, nil
}

// GetUsersIdSegments implements StrictServerInterface.
func (a *API) GetUsersIdSegments(ctx context.Context, request GetUsersIdSegmentsRequestObject) (GetUsersIdSegmentsResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return GetUsersIdSegments404Response{}, nil
	}

	segments, err := a.assignedSegments(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return GetUsersIdSegments200JSONResponse(segments), nil
}

// PostSegmentsSlug implements StrictServerInterface.
func (a *API) PostSegmentsSlug(ctx context.Context, request PostSegmentsSlugRequestObject) (PostSegmentsSlugResponseObject, error) {
	exists, err := a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if exists {
		return PostSegmentsSlug409Response{}, nil
	}

	if err = a.createSegment(ctx, request.Slug, SegmentCreation(*request.Body)); err != nil {
		return nil, err
	}

	return PostSegmentsSlug201Response{}, nil
}

// PostUsersIdSegmentsSlug implements StrictServerInterface.
func (a *API) PostUsersIdSegmentsSlug(ctx context.Context, request PostUsersIdSegmentsSlugRequestObject) (PostUsersIdSegmentsSlugResponseObject, error) {
	exists, err := a.userExists(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return PostUsersIdSegmentsSlug404Response{}, nil
	}

	exists, err = a.segmentExists(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	if !exists {
		return PostUsersIdSegmentsSlug404Response{}, nil
	}

	segments, err := a.assignedSegments(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	for _, segment := range segments {
		if segment.Slug == request.Slug {
			return PostUsersIdSegmentsSlug409Response{}, nil
		}
	}

	if err = a.assignSegment(
		ctx,
		request.Id,
		request.Slug,
		SegmentAssignment(*request.Body),
	); err != nil {
		return nil, err
	}

	return PostUsersIdSegmentsSlug200Response{}, nil
}
