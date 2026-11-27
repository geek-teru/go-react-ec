package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// Pinger を満たす stub。PingContext は err をそのまま返す。
type stubPinger struct {
	err error
}

func (s stubPinger) PingContext(context.Context) error {
	return s.err
}

func TestHealthHandler(t *testing.T) {
	errDown := errors.New("connection refused")

	tests := []struct {
		name       string
		call       func(h *HealthHandler, c echo.Context) error
		pingErr    error
		wantStatus int
		wantBody   HealthResponse
	}{
		{
			name:       "Health: DB が生きていれば 200 ok",
			call:       (*HealthHandler).Health,
			pingErr:    nil,
			wantStatus: http.StatusOK,
			wantBody:   HealthResponse{Status: "ok"},
		},
		{
			name:       "Health: DB が落ちていても 200 ok",
			call:       (*HealthHandler).Health,
			pingErr:    errDown,
			wantStatus: http.StatusOK,
			wantBody:   HealthResponse{Status: "ok"},
		},
		{
			name:       "DB: ping 成功なら 200 ok",
			call:       (*HealthHandler).DB,
			pingErr:    nil,
			wantStatus: http.StatusOK,
			wantBody:   HealthResponse{Status: "ok"},
		},
		{
			name:       "DB: ping 失敗なら 503 unavailable",
			call:       (*HealthHandler).DB,
			pingErr:    errDown,
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   HealthResponse{Status: "unavailable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ログはテストの出力に混ぜない
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h := NewHealthHandler(stubPinger{err: tt.pingErr}, logger)

			// パスはルーター層の責務なので、ハンドラーを直接呼ぶ
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := tt.call(h, c); err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var got HealthResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("body is not JSON: %v (body=%q)", err, rec.Body.String())
			}
			if got != tt.wantBody {
				t.Errorf("body = %+v, want %+v", got, tt.wantBody)
			}
		})
	}
}
