package server

import (
	"context"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server/api"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, addr string, options Options) error {
	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = nil

	e := echo.New()
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			type Error struct {
				Error string
			}

			return c.JSON(err.Code, Error{
				Error: "internal server error",
			})
		},
	}))
	apiHandler := api.NewStrictHandler(New(options), nil)
	api.RegisterHandlers(e, apiHandler)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Logger.Info().Str("addr", addr).Msg("server is up and running")
		return e.Start(addr)
	})
	g.Go(func() error {
		<-gCtx.Done()

		log.Logger.Info().Str("db", "postgres").Msg("closing connection")
		options.Connections.Postgres.Close(context.Background())

		log.Logger.Info().Msg("shutting down the server")
		return e.Shutdown(context.Background())
	})

	return g.Wait()
}
