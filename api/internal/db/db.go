package db

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	User     string `env:"DB_USER" envDefault:"local"`
	Password string `env:"DB_PASSWORD" envDefault:"local"`
	Addr     string `env:"DB_ADDR" envDefault:"127.0.0.1"`
	DBPort   uint16 `env:"DB_PORT" envDefault:"5432"`
	DBName   string `env:"DB_NAME" envDefault:"local"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

// 標準の database/sql の接続を返す（goose のマイグレーションでも使える）
func ConnectDB() (*sql.DB, error) {
	c, err := env.ParseAs[DbConfig]()
	if err != nil {
		return nil, err
	}

	cfg, err := pgx.ParseConfig("sslmode=" + c.SSLMode) // ConnConfig は ParseConfig で作る決まり
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	cfg.Host = c.Addr
	cfg.Port = c.DBPort
	cfg.User = c.User
	cfg.Password = c.Password
	cfg.Database = c.DBName
	cfg.RuntimeParams["timezone"] = "Asia/Tokyo" // MySQL の Loc にあたる

	return stdlib.OpenDB(*cfg), nil
}

// GORM で使う接続を返す
func ConnectGORM() (*gorm.DB, error) {
	sqlDB, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
}
