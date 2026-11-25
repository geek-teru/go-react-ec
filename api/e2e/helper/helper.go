//go:build e2e

package helper

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/jackc/pgx/v5"
)

// Setup は各テストの先頭で呼ぶ。テーブルを空にして、httpexpect を返す
func Setup(t *testing.T, db *sql.DB, baseURL string) *httpexpect.Expect {
	t.Helper()
	Truncate(t, db)
	return httpexpect.Default(t, baseURL)
}

// Truncate は goose の管理テーブル以外を TRUNCATE ... RESTART IDENTITY CASCADE する
func Truncate(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public' AND tablename <> 'goose_db_version'`)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		tables = append(tables, pgx.Identifier{"public", name}.Sanitize()) // 識別子をクオートする
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	if len(tables) == 0 { // 今のマイグレーションはテーブルを作らないので、空のことがある
		return
	}

	q := "TRUNCATE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE"
	if _, err := db.ExecContext(ctx, q); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}
