-include .env
export

MIGRATIONS_DIR := api/db/migrations

.PHONY: up down api-dev web-dev migrate-up migrate-down migrate-status test e2e lint fmt

up: ## Postgres と stripe-mock を起動
	docker compose up -d --wait

down: ## 停止（データは残す。消すなら docker compose down -v）
	docker compose down

api-dev: ## API をホットリロードで起動（air）
	cd api && go tool air

web-dev: ## UI を起動（後続タスクで実装）
	@echo "TODO: cd web && npm run dev"

migrate-up: ## マイグレーションを適用
	cd api && go run ./cmd/migrate up

migrate-down: ## 1つ戻す
	cd api && go run ./cmd/migrate down

migrate-status: ## 適用状況を見る
	cd api && go run ./cmd/migrate status

lint: ## 静的解析（golangci-lint）
	cd api && golangci-lint run ./...

fmt: ## 整形（gofmt + goimports）
	cd api && golangci-lint fmt ./...

test: ## ユニットテスト（後続タスクで実装）
	@echo "TODO: cd api && go test ./... && cd ../web && npm test"

e2e: ## API の E2E（Docker が必要。testcontainers がコンテナを起動する）
	cd api && go test -tags=e2e -count=1 ./e2e/...