package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server/api"
	"github.com/metafates/avito-task/swagger"
	"golang.org/x/sync/errgroup"
)

// Run the server with the given addr and options.
// It will gracefully shutdown the server when the context is done.
func Run(ctx context.Context, addr string, options Options) error {
	swaggerSpec, err := api.GetSwagger()
	if err != nil {
		return err
	}
	swaggerSpec.Servers = nil

	e := echo.New()
	apiGroup := e.Group("")
	apiGroup.Use(middleware.OapiRequestValidatorWithOptions(swaggerSpec, &middleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			var msg string
			if err.Code == http.StatusInternalServerError {
				log.Logger.Err(err).Send()

				// do not expose internal error message
				// as it can contain sensible data
				msg = "internal server error"
			} else {
				msg = fmt.Sprint(err.Message)
			}

			return c.JSON(err.Code, api.Error{
				Error: msg,
			})
		},
	}))

	apiHandler := api.NewStrictHandler(
		New(options),
		[]api.StrictMiddlewareFunc{
			middlewareLogger,
		},
	)

	api.RegisterHandlers(apiGroup, apiHandler)

	swaggerHandler, err := swagger.NewHandler(swaggerSpec)
	if err != nil {
		return err
	}

	e.GET("/docs/*", echo.WrapHandler(http.StripPrefix("/docs", swaggerHandler)))

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Logger.Info().Str("addr", addr).Msg("server is up and running")
		return e.Start(addr)
	})
	g.Go(func() error {
		<-gCtx.Done()

		log.Logger.Info().Msg("shutting down the server")
		return e.Shutdown(context.Background())
	})

	return g.Wait()
}
