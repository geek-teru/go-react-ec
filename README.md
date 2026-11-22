# go-react-ec

自己学習用の EC サイト。Go（Echo + GORM）の API と React Router v7 の UI のモノレポ。

## 技術スタック

### バックエンド（api/）

| 区分                 | 採用                                                                  |
| -------------------- | --------------------------------------------------------------------- |
| 言語・フレームワーク | Go 1.27 / Echo / GORM                                                 |
| マイグレーション     | goose（SQL ファイル）                                                 |
| 認証                 | メール + パスワード（bcrypt）、DB に保存する Cookie セッション        |
| 決済                 | Stripe テストモード（ローカル・E2E: stripe-mock）                     |
| テスト               | testing + testify / testcontainers-go                                 |
| Lint・整形           | golangci-lint（標準 + errorlint / gosec / revive、gofmt / goimports） |

### フロントエンド（web/）

| 区分           | 採用                                                         |
| -------------- | ------------------------------------------------------------ |
| フレームワーク | React Router v7（framework mode）/ TypeScript / Tailwind CSS |
| テスト         | Vitest + Testing Library                                     |

### インフラ

| 区分     | 採用                                                      |
| -------- | --------------------------------------------------------- |
| DB       | PostgreSQL（本番: Neon、ローカル・テスト: postgres:17）   |
| 実行基盤 | UI: Cloudflare Workers（SSR）、API: Cloudflare Containers |
| CI       | GitHub Actions                                            |

### 開発ツール

| 区分           | 採用                                      |
| -------------- | ----------------------------------------- |
| バージョン管理 | mise（Go / Node / goose / golangci-lint） |
| ローカル環境   | Docker Compose（db / stripe-mock）        |

UI の Worker が API の Container を呼び、`/api/*` を中継する BFF 構成です。ブラウザは Worker とだけ通信するので、CORS の設定は要りません。

## 構成

```
.
├── api/          Go API（Echo + GORM、goose）
├── web/          UI（React Router v7 + Tailwind）
├── compose.yaml  db / stripe-mock
├── Makefile      開発コマンドの入口
└── mise.toml     ツールのバージョン
```

## 必要なもの

- [mise](https://mise.jdx.dev/)
- Docker（Compose v2）

## 起動手順

```bash
mise install            # Go / Node / goose / golangci-lint
cp .env.example .env
make up                 # db と stripe-mock
make migrate-up
make api-dev            # 別ターミナル
make web-dev            # 別ターミナル
```

## Make コマンド

| コマンド                        | 内容                          |
| ------------------------------- | ----------------------------- |
| `make up` / `make down`         | db・stripe-mock の起動 / 停止 |
| `make api-dev` / `make web-dev` | API / UI の開発サーバー       |
| `make migrate-up`               | マイグレーションの適用        |
| `make test` / `make e2e`        | テスト / API の E2E           |
| `make lint` / `make fmt`        | 静的解析 / 整形               |
