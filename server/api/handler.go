package api

import (
	"encoding/json"
	"net/http"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/metafates/avito-task/log"
)

func NewHandler(options Options) (http.Handler, error) {
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil

	api := &API{options: options}

	strictHandler := NewStrictHandlerWithOptions(
		api,
		[]StrictMiddlewareFunc{
			// middlewareLogger,
		},
		StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				// do not expose internal error messages to the client since
				// they *can* contain sensible data
				log.Logger.Err(err).Msg("error")
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
	)

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			type Error struct {
				Error string
			}

			err := json.NewEncoder(w).Encode(Error{
				Error: message,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
	}))

	return HandlerFromMux(strictHandler, r), nil
}
