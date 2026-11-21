-include .env
export

MIGRATIONS_DIR := api/db/migrations

.PHONY: up down api-dev web-dev migrate-up test e2e

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

test: ## ユニットテスト（後続タスクで実装）
	@echo "TODO: cd api && go test ./... && cd ../web && npm test"

e2e: ## API の E2E（後続タスクで実装）
	@echo "TODO: cd api && go test ./e2e/..."