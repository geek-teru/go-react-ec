package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// DB の疎通確認に使うインターフェース。
// *sql.DB がそのまま満たす。テストではモックを渡す。
type Pinger interface {
	PingContext(ctx context.Context) error
}

type HealthResponse struct {
	Status string `json:"status"`
}

type HealthHandler struct {
	db      Pinger
	logger  *slog.Logger
	timeout time.Duration
}

func NewHealthHandler(db Pinger, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{db: db, logger: logger, timeout: 2 * time.Second}
}

// healthチェック：プロセスが生きていれば 200（DB は見ない）
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// DB healthチェック：DB に ping して、成功なら 200、失敗なら 503
func (h *HealthHandler) DB(c echo.Context) error {
	// DB が応答しないときにリクエストが詰まらないよう、上限を切る
	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		h.logger.WarnContext(ctx, "db ping failed", slog.Any("err", err))
		return c.JSON(http.StatusServiceUnavailable, HealthResponse{Status: "unavailable"})
	}
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
