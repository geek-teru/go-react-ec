package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geek-teru/go-react-ec/api/internal/config"
	"github.com/geek-teru/go-react-ec/api/internal/db"
	"github.com/geek-teru/go-react-ec/api/internal/httpserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// run server
	if err := run(logger); err != nil {
		logger.Error("server stopped", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	if _, err := db.ConnectGORM(); err != nil {
		return err
	} // 後続で各機能に渡す

	e := httpserver.New(logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("start failed", slog.Any("err", err))
			stop()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return e.Shutdown(shutdownCtx)
}
