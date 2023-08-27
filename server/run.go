package server

import (
	"context"
	"net/http"

	"github.com/metafates/avito-task/log"
	"github.com/metafates/avito-task/server/api"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, addr string, options api.Options) error {
	handler, err := api.NewHandler(options)
	if err != nil {
		return err
	}
	server := &http.Server{Addr: addr, Handler: handler}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Logger.Info().Str("addr", addr).Msg("server is up and running")
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()

		log.Logger.Info().Str("db", "postgres").Msg("closing connection")
		options.Connections.Postgres.Close(context.Background())

		log.Logger.Info().Msg("shutting down the server")
		return server.Shutdown(context.Background())
	})

	return g.Wait()
}
