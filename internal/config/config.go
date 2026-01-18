package config

import "os"

type Config struct {
	Addr string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	
	cfg := &Config{
		Addr: port,
	}
	
	return cfg
}