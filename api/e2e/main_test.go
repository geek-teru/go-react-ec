//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/geek-teru/go-react-ec/api/internal/httpserver"
	"github.com/geek-teru/go-react-ec/api/internal/testutil"
)

// パッケージ内のテストで共有する
var (
	sqlDB   *sql.DB
	baseURL string
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	pg, err := testutil.StartPostgres(ctx)
	if err != nil {
		log.Printf("setup: %v", err)
		return 1
	}
	defer func() {
		if err := pg.Close(); err != nil {
			log.Printf("teardown: %v", err)
		}
	}()

	if err := testutil.Migrate(pg.DB); err != nil {
		log.Printf("migrate: %v", err)
		return 1
	}
	sqlDB = pg.DB

	// cmd/server と同じ組み立て関数を通す。ログはテストの出力を汚さないよう捨てる
	e := httpserver.New(slog.New(slog.DiscardHandler), pg.DB)
	srv := httptest.NewServer(e)
	defer srv.Close()
	baseURL = srv.URL

	return m.Run()
}
