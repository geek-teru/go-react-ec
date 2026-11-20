# Tips

- [go mod init と依存の追加](#go-mod-init-と依存の追加)
- [air とは](#air-とは)
- [go get の `-tool` あり・なし](#go-get-の--tool-ありなし)
- [slog とは](#slog-とは)

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

- 設定は `api/.air.toml`（見張る拡張子、ビルドのコマンド、起動するバイナリ）
- 起動は `go tool air`（`make api-dev`）
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
