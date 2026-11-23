package config

import "github.com/caarlos0/env/v11"

type Config struct {
	// 共通
	ServiceEnv string `env:"SERVICE_ENV" envDefault:"local"`
	Port       string `env:"API_PORT" envDefault:"8080"`

	// Stripe 関連
	// StripeSecretKey     string `env:"STRIPE_SECRET_KEY,required"`
	// StripeWebhookSecret string `env:"STRIPE_WEBHOOK_SECRET"`
	// StripeAPIBase       string `env:"STRIPE_API_BASE"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
