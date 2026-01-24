package config

import "os"

type Config struct {
	Addr       string
	AuthConfig *AuthConfig
}

type AuthConfig struct {
	SecretKey string
	PublicKey string
	IssuerUrl string
}

func LoadConfig(authConfig *AuthConfig) *Config {
	port := os.Getenv("PORT")

	cfg := &Config{
		Addr: port,
		AuthConfig: authConfig,
	}

	return cfg
}
