# Tips

- [go mod init と依存の追加](#go-mod-init-と依存の追加)
- [air とは](#air-とは)
- [go get の `-tool` あり・なし](#go-get-の--tool-ありなし)
- [slog とは](#slog-とは)
- [goose とは](#goose-とは)

## go mod init と依存の追加

```bash
cd api
go mod init github.com/geek-teru/go-react-ec/api   # go.mod を作る（最初の1回だけ）
go get github.com/labstack/echo/v4                 # 依存を追加する（go.mod / go.sum に記録される）
```

- `go mod init` が作るのは go.mod だけ。依存は `go get` で入れる
- 引数はモジュールパス。import のときの接頭辞になる（例: `github.com/geek-teru/go-react-ec/api/internal/config`）

## air とは

Go のホットリロードツール。`.go` ファイルを保存すると、自動でビルドしてサーバーを再起動する。

```bash
cd api && go get -tool github.com/air-verse/air@latest   # go.mod の tool に入る（最初の1回だけ）
make api-dev                                              # = cd api && go tool air
```

- 設定は `api/.air.toml`。`./cmd/server` を `api/tmp/server` にビルドして起動する。`tmp/` と `db/migrations` は見張らない
- 保存すると `xxx.go has changed` → `building...` → `running...` とログが出る。ビルドに失敗したら前のサーバーは止まる（`stop_on_error`）
- WSL から `/mnt/c` のファイルを見ると、Windows 側（VS Code など）の編集が inotify で届かない。そのため `poll = true` でポーリングにしている
- 再起動のときは SIGINT を送る（`send_interrupt`）ので、`main.go` の graceful shutdown が効く
- `api/tmp/` は `.gitignore` 済み。終了時に消える（`clean_on_exit`）
- 開発中のローカル専用。本番のコンテナには入れない

## go get の `-tool` あり・なし

| | `go get`（なし） | `go get -tool`（あり、Go 1.24〜） |
| --- | --- | --- |
| 対象 | コードから import するライブラリ | 開発中にコマンドとして実行するツール |
| go.mod | `require` に入る | `tool` に入る |
| 使い方 | `import "..."` | `go tool <名前>` |
| アプリのバイナリ | 含まれる | 含まれない |
| 例 | echo、gorm、goose、pgx | air |

- `-tool` を使えば、`go install` しなくても、clone した人全員が同じバージョンのツールを使える
- goose は CLI もあるが、このリポジトリでは `cmd/migrate` から呼ぶのでライブラリとして入れる

## slog とは

Go 1.21 から標準ライブラリに入った構造化ロガー（`log/slog`）。外部の依存なしで、キーと値の組でログを出せる。

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("request", slog.String("method", "GET"), slog.Int("status", 200))
// {"time":"...","level":"INFO","msg":"request","method":"GET","status":200}
```

- 新規なら slog が第一候補。以前の主流は zap・zerolog（速度重視の場面では今も使われる）
- 出力の形は Handler で決まる（`NewJSONHandler` / `NewTextHandler`）。zap などを Handler として差し替えても、呼び出す側は変わらない
- このリポジトリでは `cmd/server` で JSON のロガーを作り、`httpserver.New(logger)` に渡してアクセスログとエラーログに使う
- `logger.Info("msg", "key")` のように値の無いキーを渡すと `!BADKEY` になる。`slog.String` などの型付きの関数を使うと防げる

## goose とは

Go 製の DB マイグレーションツール。SQL ファイルに Up と Down を書き、番号順に適用する。CLI とライブラリの両方がある。

```sql
-- api/db/migrations/00001_create_users.sql
-- +goose Up
CREATE TABLE users (id BIGSERIAL PRIMARY KEY, email TEXT NOT NULL UNIQUE);

-- +goose Down
DROP TABLE IF EXISTS users;
```

| やること | コマンド |
| --- | --- |
| 雛形を作る | `goose -dir api/db/migrations -s create <名前> sql`（`-s` で連番） |
| 適用する | `make migrate-up`（`go run ./cmd/migrate up`） |
| 1つ戻す | `make migrate-down` |
| 状態を見る | `make migrate-status` |

- このリポジトリでは、SQL を `embed` で埋め込み、`cmd/migrate` から `goose.SetBaseFS` → `goose.SetDialect("postgres")` → `goose.Up(db, ".")` で流す。接続は API と同じ `db.ConnectDB()`
- 適用済みの番号は `goose_db_version` テーブルに記録される。適用済みのファイルを直しても再適用されないので、変更は新しいファイルで足す
- 関数・トリガーなど `;` を含む文は `-- +goose StatementBegin` / `-- +goose StatementEnd` で囲む
- GORM の AutoMigrate は使わない（履歴・Down が無く、列の削除や改名ができないため）
