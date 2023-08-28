// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// CSV defines model for CSV.
type CSV = string

// Date defines model for Date.
type Date = openapi_types.Date

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// Outreach defines model for Outreach.
type Outreach = float32

// Slug defines model for Slug.
type Slug = string

// Timestamp defines model for Timestamp.
type Timestamp = time.Time

// UserID defines model for UserID.
type UserID = int32

// UserSegment defines model for UserSegment.
type UserSegment struct {
	ExpiresAt *Timestamp `json:"expiresAt,omitempty"`
	Slug      Slug       `json:"slug"`
}

// User defines model for User.
type User = UserID

// SegmentAssignment defines model for SegmentAssignment.
type SegmentAssignment struct {
	ExpiresAt *Timestamp `json:"expiresAt,omitempty"`
}

// SegmentCreation defines model for SegmentCreation.
type SegmentCreation struct {
	Outreach *Outreach `json:"outreach,omitempty"`
}

// GetAuditParams defines parameters for GetAudit.
type GetAuditParams struct {
	// From Start date of the audit window
	From *Date `form:"from,omitempty" json:"from,omitempty"`

	// To End date of the audit window
	To *Date `form:"to,omitempty" json:"to,omitempty"`

	// User Show audit for specifc user with the given id
	User *UserID `form:"user,omitempty" json:"user,omitempty"`
}

// PostSegmentsSlugJSONBody defines parameters for PostSegmentsSlug.
type PostSegmentsSlugJSONBody struct {
	Outreach *Outreach `json:"outreach,omitempty"`
}

// PostUsersIdSegmentsSlugJSONBody defines parameters for PostUsersIdSegmentsSlug.
type PostUsersIdSegmentsSlugJSONBody struct {
	ExpiresAt *Timestamp `json:"expiresAt,omitempty"`
}

// PostSegmentsSlugJSONRequestBody defines body for PostSegmentsSlug for application/json ContentType.
type PostSegmentsSlugJSONRequestBody PostSegmentsSlugJSONBody

// PostUsersIdSegmentsSlugJSONRequestBody defines body for PostUsersIdSegmentsSlug for application/json ContentType.
type PostUsersIdSegmentsSlugJSONRequestBody PostUsersIdSegmentsSlugJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get audit of changes
	// (GET /audit)
	GetAudit(ctx echo.Context, params GetAuditParams) error
	// Delete a segment
	// (DELETE /segments/{slug})
	DeleteSegmentsSlug(ctx echo.Context, slug Slug) error
	// Create a new segment
	// (POST /segments/{slug})
	PostSegmentsSlug(ctx echo.Context, slug Slug) error
	// Create a new user
	// (POST /users/{id})
	PostUsersId(ctx echo.Context, id User) error
	// Get active segments assigned to a user
	// (GET /users/{id}/segments)
	GetUsersIdSegments(ctx echo.Context, id User) error
	// Deprive segment from a user
	// (DELETE /users/{id}/segments/{slug})
	DeleteUsersIdSegmentsSlug(ctx echo.Context, id User, slug Slug) error
	// Assign segment to a user
	// (POST /users/{id}/segments/{slug})
	PostUsersIdSegmentsSlug(ctx echo.Context, id User, slug Slug) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAudit converts echo context to params.
func (w *ServerInterfaceWrapper) GetAudit(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAuditParams
	// ------------- Optional query parameter "from" -------------

	err = runtime.BindQueryParameter("form", true, false, "from", ctx.QueryParams(), &params.From)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from: %s", err))
	}

	// ------------- Optional query parameter "to" -------------

	err = runtime.BindQueryParameter("form", true, false, "to", ctx.QueryParams(), &params.To)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter to: %s", err))
	}

	// ------------- Optional query parameter "user" -------------

	err = runtime.BindQueryParameter("form", true, false, "user", ctx.QueryParams(), &params.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAudit(ctx, params)
	return err
}

// DeleteSegmentsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteSegmentsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "slug" -------------
	var slug Slug

	err = runtime.BindStyledParameterWithLocation("simple", false, "slug", runtime.ParamLocationPath, ctx.Param("slug"), &slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteSegmentsSlug(ctx, slug)
	return err
}

// PostSegmentsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) PostSegmentsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "slug" -------------
	var slug Slug

	err = runtime.BindStyledParameterWithLocation("simple", false, "slug", runtime.ParamLocationPath, ctx.Param("slug"), &slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostSegmentsSlug(ctx, slug)
	return err
}

// PostUsersId converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id User

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUsersId(ctx, id)
	return err
}

// GetUsersIdSegments converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersIdSegments(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id User

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUsersIdSegments(ctx, id)
	return err
}

// DeleteUsersIdSegmentsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUsersIdSegmentsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id User

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "slug" -------------
	var slug Slug

	err = runtime.BindStyledParameterWithLocation("simple", false, "slug", runtime.ParamLocationPath, ctx.Param("slug"), &slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteUsersIdSegmentsSlug(ctx, id, slug)
	return err
}

// PostUsersIdSegmentsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersIdSegmentsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id User

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "slug" -------------
	var slug Slug

	err = runtime.BindStyledParameterWithLocation("simple", false, "slug", runtime.ParamLocationPath, ctx.Param("slug"), &slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUsersIdSegmentsSlug(ctx, id, slug)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/audit", wrapper.GetAudit)
	router.DELETE(baseURL+"/segments/:slug", wrapper.DeleteSegmentsSlug)
	router.POST(baseURL+"/segments/:slug", wrapper.PostSegmentsSlug)
	router.POST(baseURL+"/users/:id", wrapper.PostUsersId)
	router.GET(baseURL+"/users/:id/segments", wrapper.GetUsersIdSegments)
	router.DELETE(baseURL+"/users/:id/segments/:slug", wrapper.DeleteUsersIdSegmentsSlug)
	router.POST(baseURL+"/users/:id/segments/:slug", wrapper.PostUsersIdSegmentsSlug)

}

type GetAuditRequestObject struct {
	Params GetAuditParams
}

type GetAuditResponseObject interface {
	VisitGetAuditResponse(w http.ResponseWriter) error
}

type GetAudit200TextcsvResponse struct {
	Body          io.Reader
	ContentLength int64
}

func (response GetAudit200TextcsvResponse) VisitGetAuditResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/csv")
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type DeleteSegmentsSlugRequestObject struct {
	Slug Slug `json:"slug"`
}

type DeleteSegmentsSlugResponseObject interface {
	VisitDeleteSegmentsSlugResponse(w http.ResponseWriter) error
}

type DeleteSegmentsSlug200Response struct {
}

func (response DeleteSegmentsSlug200Response) VisitDeleteSegmentsSlugResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type DeleteSegmentsSlug404JSONResponse Error

func (response DeleteSegmentsSlug404JSONResponse) VisitDeleteSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type PostSegmentsSlugRequestObject struct {
	Slug Slug `json:"slug"`
	Body *PostSegmentsSlugJSONRequestBody
}

type PostSegmentsSlugResponseObject interface {
	VisitPostSegmentsSlugResponse(w http.ResponseWriter) error
}

type PostSegmentsSlug201Response struct {
}

func (response PostSegmentsSlug201Response) VisitPostSegmentsSlugResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type PostSegmentsSlug409JSONResponse Error

func (response PostSegmentsSlug409JSONResponse) VisitPostSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersIdRequestObject struct {
	Id User `json:"id"`
}

type PostUsersIdResponseObject interface {
	VisitPostUsersIdResponse(w http.ResponseWriter) error
}

type PostUsersId201Response struct {
}

func (response PostUsersId201Response) VisitPostUsersIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type PostUsersId409JSONResponse Error

func (response PostUsersId409JSONResponse) VisitPostUsersIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type GetUsersIdSegmentsRequestObject struct {
	Id User `json:"id"`
}

type GetUsersIdSegmentsResponseObject interface {
	VisitGetUsersIdSegmentsResponse(w http.ResponseWriter) error
}

type GetUsersIdSegments200JSONResponse []UserSegment

func (response GetUsersIdSegments200JSONResponse) VisitGetUsersIdSegmentsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUsersIdSegments404JSONResponse Error

func (response GetUsersIdSegments404JSONResponse) VisitGetUsersIdSegmentsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type DeleteUsersIdSegmentsSlugRequestObject struct {
	Id   User `json:"id"`
	Slug Slug `json:"slug"`
}

type DeleteUsersIdSegmentsSlugResponseObject interface {
	VisitDeleteUsersIdSegmentsSlugResponse(w http.ResponseWriter) error
}

type DeleteUsersIdSegmentsSlug200Response struct {
}

func (response DeleteUsersIdSegmentsSlug200Response) VisitDeleteUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type DeleteUsersIdSegmentsSlug404JSONResponse Error

func (response DeleteUsersIdSegmentsSlug404JSONResponse) VisitDeleteUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersIdSegmentsSlugRequestObject struct {
	Id   User `json:"id"`
	Slug Slug `json:"slug"`
	Body *PostUsersIdSegmentsSlugJSONRequestBody
}

type PostUsersIdSegmentsSlugResponseObject interface {
	VisitPostUsersIdSegmentsSlugResponse(w http.ResponseWriter) error
}

type PostUsersIdSegmentsSlug200Response struct {
}

func (response PostUsersIdSegmentsSlug200Response) VisitPostUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PostUsersIdSegmentsSlug400JSONResponse Error

func (response PostUsersIdSegmentsSlug400JSONResponse) VisitPostUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersIdSegmentsSlug404JSONResponse Error

func (response PostUsersIdSegmentsSlug404JSONResponse) VisitPostUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersIdSegmentsSlug409JSONResponse Error

func (response PostUsersIdSegmentsSlug409JSONResponse) VisitPostUsersIdSegmentsSlugResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get audit of changes
	// (GET /audit)
	GetAudit(ctx context.Context, request GetAuditRequestObject) (GetAuditResponseObject, error)
	// Delete a segment
	// (DELETE /segments/{slug})
	DeleteSegmentsSlug(ctx context.Context, request DeleteSegmentsSlugRequestObject) (DeleteSegmentsSlugResponseObject, error)
	// Create a new segment
	// (POST /segments/{slug})
	PostSegmentsSlug(ctx context.Context, request PostSegmentsSlugRequestObject) (PostSegmentsSlugResponseObject, error)
	// Create a new user
	// (POST /users/{id})
	PostUsersId(ctx context.Context, request PostUsersIdRequestObject) (PostUsersIdResponseObject, error)
	// Get active segments assigned to a user
	// (GET /users/{id}/segments)
	GetUsersIdSegments(ctx context.Context, request GetUsersIdSegmentsRequestObject) (GetUsersIdSegmentsResponseObject, error)
	// Deprive segment from a user
	// (DELETE /users/{id}/segments/{slug})
	DeleteUsersIdSegmentsSlug(ctx context.Context, request DeleteUsersIdSegmentsSlugRequestObject) (DeleteUsersIdSegmentsSlugResponseObject, error)
	// Assign segment to a user
	// (POST /users/{id}/segments/{slug})
	PostUsersIdSegmentsSlug(ctx context.Context, request PostUsersIdSegmentsSlugRequestObject) (PostUsersIdSegmentsSlugResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictEchoHandlerFunc
type StrictMiddlewareFunc = runtime.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetAudit operation middleware
func (sh *strictHandler) GetAudit(ctx echo.Context, params GetAuditParams) error {
	var request GetAuditRequestObject

	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetAudit(ctx.Request().Context(), request.(GetAuditRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAudit")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetAuditResponseObject); ok {
		return validResponse.VisitGetAuditResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// DeleteSegmentsSlug operation middleware
func (sh *strictHandler) DeleteSegmentsSlug(ctx echo.Context, slug Slug) error {
	var request DeleteSegmentsSlugRequestObject

	request.Slug = slug

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteSegmentsSlug(ctx.Request().Context(), request.(DeleteSegmentsSlugRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteSegmentsSlug")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(DeleteSegmentsSlugResponseObject); ok {
		return validResponse.VisitDeleteSegmentsSlugResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostSegmentsSlug operation middleware
func (sh *strictHandler) PostSegmentsSlug(ctx echo.Context, slug Slug) error {
	var request PostSegmentsSlugRequestObject

	request.Slug = slug

	var body PostSegmentsSlugJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostSegmentsSlug(ctx.Request().Context(), request.(PostSegmentsSlugRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostSegmentsSlug")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostSegmentsSlugResponseObject); ok {
		return validResponse.VisitPostSegmentsSlugResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostUsersId operation middleware
func (sh *strictHandler) PostUsersId(ctx echo.Context, id User) error {
	var request PostUsersIdRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostUsersId(ctx.Request().Context(), request.(PostUsersIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostUsersId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostUsersIdResponseObject); ok {
		return validResponse.VisitPostUsersIdResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetUsersIdSegments operation middleware
func (sh *strictHandler) GetUsersIdSegments(ctx echo.Context, id User) error {
	var request GetUsersIdSegmentsRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetUsersIdSegments(ctx.Request().Context(), request.(GetUsersIdSegmentsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUsersIdSegments")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetUsersIdSegmentsResponseObject); ok {
		return validResponse.VisitGetUsersIdSegmentsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// DeleteUsersIdSegmentsSlug operation middleware
func (sh *strictHandler) DeleteUsersIdSegmentsSlug(ctx echo.Context, id User, slug Slug) error {
	var request DeleteUsersIdSegmentsSlugRequestObject

	request.Id = id
	request.Slug = slug

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteUsersIdSegmentsSlug(ctx.Request().Context(), request.(DeleteUsersIdSegmentsSlugRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteUsersIdSegmentsSlug")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(DeleteUsersIdSegmentsSlugResponseObject); ok {
		return validResponse.VisitDeleteUsersIdSegmentsSlugResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostUsersIdSegmentsSlug operation middleware
func (sh *strictHandler) PostUsersIdSegmentsSlug(ctx echo.Context, id User, slug Slug) error {
	var request PostUsersIdSegmentsSlugRequestObject

	request.Id = id
	request.Slug = slug

	var body PostUsersIdSegmentsSlugJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostUsersIdSegmentsSlug(ctx.Request().Context(), request.(PostUsersIdSegmentsSlugRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostUsersIdSegmentsSlug")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostUsersIdSegmentsSlugResponseObject); ok {
		return validResponse.VisitPostUsersIdSegmentsSlugResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xXTXPbNhD9Kxi0R8aUk1zKU53I0/FMO+lUaS8ZH2BySSIhAAZYSlY9/O+dBUBLoqgP",
	"O6knJ0nQEvt23+7b5QPPjWqNBo2OZw+8FVYoQLD+16LpKvoswOVWtiiN5hlfQKVAI3P0b8IlnbUCa55w",
	"LRTwjMd/LHztpIWCZ2g7SLjLa1CCLvzZQskz/lO68Z6Gf13qnfZ9wv92YPe90ym7mTM0bGXsF7aS3vME",
	"Clk8GwM5uZnznmDQDeDwnSkkhKyE+K+ck5Wmb3SYG43xq2jbRuaC8KafHYF+2HLbWtOCxXgX3LfSgrvC",
	"U4g+SgUOhWp9anDdUoTm7jPkGGDuZimAYy5ShYYJ1lE++2TA/96CCNbPRm86tCDy+hT4D4PdWdg9MGCC",
	"aVgNEfgn43Xk7f3iH/oojVUCecZzt+SPVzu0UlcU6lwg7NgVdDBheG2tsRP0DMejB/rtwvoUzW73Ykv4",
	"h60MPYIoGyOQJ1yJe6k6xbPLhCupw/fZ4y26U3eRsdiHSurfQVdY+0f2gtjUyDjkVyjVZNyx0LftpcY3",
	"rze2UiNUAQYZx+L5PpWcBKU4VxC2U+4fvJ2sJrhHsFo0c5N7ZJ1teMZrxNZlaVpJrLu7i9yoVAGKUiC4",
	"VCwlmlco3Jf0rjF3qRJSp39dX83/uL5QBQGVujRDp4jcBwlKSLp5ufz3ooDlrxX9povJHiU2BOyKbmYf",
	"wRHhS7Au1Pjs4vJiRnamBS1ayTP+xh8lXsQ87lR0hfSeKvAflG7fmDcFz/hvgFfeINlR7U97co3CIqMq",
	"YKZkWAPz97KV1IVZDdL5tQO73mhnaY3i56ql77K+T8aer3XxVL9ovtnrojar6Ko0lrkWclnmXv38sPBQ",
	"KrkEzfyEmILhpfLJw+KWStS1RrvQEq9ns5G4ItxjSlq1I6rHLiehm9J3H58pWV4LXYEL+tgpJew61EZM",
	"wcgkjXrq0gfqoD7M1waCTO4W2Nyfx4Z3i2Gm7wd4YD3o8hycK7umWbPgw3fS29nbJ42cY9kJsj2RnwHE",
	"iHCKma2EY9pQdXS6GCUuxMzEZu7sddcUno1JVKvbhLfGTXTtn8bhXkqHBWN9KN6dHSQdD/B+j5bLU7QE",
	"VnI/aSMrv7w4K9IFQuBeOnRMNBZEsR4xcnAbSKlJXfogi368uZ5mye+Wx1kiE3dT8HOS65fSl86mdxpy",
	"dyxnYe3bTdijDBybLjEBQ7Xyk+p2PEaJoNw5WrrYtF6c7sJasT5SUI4Jv+9CsbvpvojUeBoO6YkX4hzl",
	"EobaPYT1WfV7gNWzxX1E8ZM1voDWyuUL6vq1xBpsGOY03LdeRU/ougf6+ADtN9+Y/OR7jIJDBDxvImy9",
	"kvZPoXGoyEDj7P+n8Z0oWIzgRyydFx6HcfBt0bBTu4ff5fu+/y8AAP//+YvSlcERAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
