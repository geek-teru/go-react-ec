package main

import (
	"log"
	"os"

	"github.com/pressly/goose/v3"

	"github.com/geek-teru/go-react-ec/api/db/migrations"
	"github.com/geek-teru/go-react-ec/api/internal/db"
)

func main() {
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	sqlDB, err := db.ConnectDB() // API と同じ DB_* の設定で繋ぐ
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrations.FS) // embed した SQL を読む
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "up":
		err = goose.Up(sqlDB, ".")
	case "down":
		err = goose.Down(sqlDB, ".") // 1つだけ戻す
	case "status":
		err = goose.Status(sqlDB, ".")
	default:
		log.Fatalf("unknown command: %q (up | down | status)", cmd) //nolint:gosec // 開発者が実行する CLI の引数。%q で改行もエスケープ済み
	}
	if err != nil {
		log.Fatal(err)
	}
}
