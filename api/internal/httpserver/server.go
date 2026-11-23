package httpserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/geek-teru/go-react-ec/api/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// errorBody(Json形式)
type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func New(logger *slog.Logger, db handler.Pinger) *echo.Echo {
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
			logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "request",
				slog.String("method", v.Method), slog.String("uri", v.URI),
				slog.Int("status", v.Status), slog.Duration("latency", v.Latency),
				slog.String("request_id", v.RequestID))
			return nil
		},
	}))
	e.Use(middleware.Recover())

	api := e.Group("/api")

	healthHandler := handler.NewHealthHandler(db, logger)

	// healthcheck
	api.GET("/health", healthHandler.Health)
	api.GET("/health/db", healthHandler.DB)

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
