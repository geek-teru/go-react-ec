# Tips

- [go mod init と依存の追加](#go-mod-init-と依存の追加)
- [air とは](#air-とは)
- [go get の `-tool` あり・なし](#go-get-の--tool-ありなし)
- [slog とは](#slog-とは)
- [goose とは](#goose-とは)
- [golangci-lint とは](#golangci-lint-とは)
- [golangci-lint で使う linter・formatter](#golangci-lint-で使う-linterformatter)

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

## golangci-lint とは

Go の linter と formatter をまとめて実行するツール。100 を超える linter を内蔵し、使うものを `api/.golangci.yml` で選ぶ。

```bash
make lint   # = cd api && golangci-lint run ./...   検査
make fmt    # = cd api && golangci-lint fmt ./...   整形（gofmt + goimports）
```

- バージョンは `mise.toml` で固定（2.6.0）。v2 は設定ファイルに `version: "2"` が要り、v1 の書き方は読めない（`golangci-lint migrate` で変換できる）
- 設定ファイルは実行したディレクトリから上へ探す。`api/` で実行するので `api/.golangci.yml` が読まれる。無くても標準のセットで動くが、formatter は使われない
- 手元でも CI でも同じ設定ファイルを読むので、誰がどこで実行しても同じルールで検査できる
- golangci-lint をビルドした Go より新しい Go を go.mod で指定すると `the Go language version ... is lower than the targeted Go version` で止まる。golangci-lint の版を上げる
- `exclusions.presets` で、よくある誤検知（`defer f.Close()` の戻り値など）や、公開関数のコメントの書式の指摘をまとめて外せる

## golangci-lint で使う linter・formatter

| 名前 | 種類 | 見るもの・整えるもの |
| --- | --- | --- |
| errcheck | 標準 | 返ってきたエラーを確かめずに捨てていないか |
| govet | 標準 | 書式の指定ミス、構造体のコピーなど（`go vet` と同じ） |
| ineffassign | 標準 | 代入した値を使わずに上書きしていないか |
| staticcheck | 標準 | 非推奨の API、意味の無いコード、よくある書き間違い（stylecheck・gosimple を含む） |
| unused | 標準 | 使っていない関数・変数・型・定数 |
| errorlint | 追加 | エラーを `==` や型アサーションで比べていないか（`errors.Is` / `errors.As` を使う） |
| gosec | 追加 | SQL の文字列連結、弱い乱数、タイムアウトの無い HTTP サーバーなど。`_test.go` では外す |
| revive | 追加 | 命名（`DbConfig` → `DBConfig`）、使っていない引数など、スタイル全般 |
| gofmt | formatter | インデント・空白・フィールドの揃えを標準の書式にする |
| goimports | formatter | import の追加・削除と並び替え。標準 → 外部 → このリポジトリ（`local-prefixes`）の順に分ける |

- 「標準」は `linters.default: standard` で動くもの。「追加」は `linters.enable` で足したもの
- 使っていない引数は `_` にすると revive の指摘が消える（例: `func(_ echo.Context, v ...)`）
