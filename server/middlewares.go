package server

import (
	"net"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/metafates/avito-task/log"
)

// middlewareLogger Logs all requests
func middlewareLogger(f runtime.StrictEchoHandlerFunc, _ string) runtime.StrictEchoHandlerFunc {
	return func(ctx echo.Context, request any) (response any, err error) {
		event := log.
			Logger.
			Info().
			Str("path", ctx.Path()).
			Str("method", ctx.Request().Method)

		if ip := net.ParseIP(ctx.RealIP()); ip != nil {
			event = event.IPAddr("ip", ip)
		}

		event.Send()

		return f(ctx, request)
	}
}
