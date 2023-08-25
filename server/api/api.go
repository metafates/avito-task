package api

import (
	"context"
)

var _ StrictServerInterface = (*API)(nil)

type API struct {
	options Options
}

// PostUsersId implements StrictServerInterface.
func (*API) PostUsersId(ctx context.Context, request PostUsersIdRequestObject) (PostUsersIdResponseObject, error) {
	panic("unimplemented")
}

// DeleteSegmentsSlug implements StrictServerInterface.
func (*API) DeleteSegmentsSlug(ctx context.Context, request DeleteSegmentsSlugRequestObject) (DeleteSegmentsSlugResponseObject, error) {
	panic("unimplemented")
}

// DeleteUsersIdSegmentsSlug implements StrictServerInterface.
func (*API) DeleteUsersIdSegmentsSlug(ctx context.Context, request DeleteUsersIdSegmentsSlugRequestObject) (DeleteUsersIdSegmentsSlugResponseObject, error) {
	panic("unimplemented")
}

// GetUsersIdSegments implements StrictServerInterface.
func (*API) GetUsersIdSegments(ctx context.Context, request GetUsersIdSegmentsRequestObject) (GetUsersIdSegmentsResponseObject, error) {
	panic("unimplemented")
}

// PostSegmentsSlug implements StrictServerInterface.
func (*API) PostSegmentsSlug(ctx context.Context, request PostSegmentsSlugRequestObject) (PostSegmentsSlugResponseObject, error) {
	panic("unimplemented")
}

// PostUsersIdSegmentsSlug implements StrictServerInterface.
func (*API) PostUsersIdSegmentsSlug(ctx context.Context, request PostUsersIdSegmentsSlugRequestObject) (PostUsersIdSegmentsSlugResponseObject, error) {
	panic("unimplemented")
}
