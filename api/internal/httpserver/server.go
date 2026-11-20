package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type errorBody struct {
	Error struct{ Code, Message string } `json:"error"`
}

func New(logger *slog.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// custom error handler
	e.HTTPErrorHandler = errorHandler(logger)

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogMethod:    true,
		LogLatency:   true,
		LogRequestID: true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "request",
				slog.String("method", v.Method), slog.String("uri", v.URI),
				slog.Int("status", v.Status), slog.Duration("latency", v.Latency),
				slog.String("request_id", v.RequestID))
			return nil
		},
	}))
	e.Use(middleware.Recover())

	api := e.Group("/api")
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	return e
}

func errorHandler(logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		status, code, msg := http.StatusInternalServerError, "internal_error", "internal server error"

		var he *echo.HTTPError
		if errors.As(err, &he) {
			status = he.Code
			code = http.StatusText(he.Code)
			if m, ok := he.Message.(string); ok {
				msg = m
			}
		} else {
			logger.Error("unhandled error", slog.Any("err", err))
		}
		var b errorBody
		b.Error.Code, b.Error.Message = code, msg
		_ = c.JSON(status, b)
	}
}
