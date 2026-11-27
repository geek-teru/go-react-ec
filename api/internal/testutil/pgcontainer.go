package testutil

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // "pgx" ドライバーを登録する
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/geek-teru/go-react-ec/api/db/migrations"
)

// 本番(Neon)・compose と同じメジャーバージョンにそろえる
const postgresImage = "postgres:17"

type Postgres struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
}

// StartPostgres は使い捨ての Postgres を起動して、接続済みの *sql.DB を返す。
// .env の DB_* は読まない（ローカルの compose の DB を誤って TRUNCATE しないため）。
func StartPostgres(ctx context.Context) (*Postgres, error) {
	ctr, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(), // 「ready」のログ2回 + ポートの待ち受け
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres container: %w", err)
	}

	// internal/db と同じく timezone を Asia/Tokyo にする
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable", "timezone=Asia/Tokyo")
	if err != nil {
		_ = testcontainers.TerminateContainer(ctr)
		return nil, fmt.Errorf("connection string: %w", err)
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		_ = testcontainers.TerminateContainer(ctr)
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		_ = testcontainers.TerminateContainer(ctr)
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return &Postgres{Container: ctr, DB: sqlDB}, nil
}

// Close は接続を閉じてコンテナを消す
func (p *Postgres) Close() error {
	_ = p.DB.Close()
	return testcontainers.TerminateContainer(p.Container)
}

// Migrate は cmd/migrate と同じ embed の SQL で goose up する
func Migrate(sqlDB *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, ".")
}
