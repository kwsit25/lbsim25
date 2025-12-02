package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort      int `env:"ODJ_EE_HTTP_PORT" envDefault:"8080"`
	ResponseDelay int `env:"ODJ_DEP_RESPONSE_DELAY" envDefault:"10"`
}

func Get() (*Config, error) {
	var cfg Config

	// load .env file
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		// no .env file found, a vaild case for incontainerized environments.
		//nolint:forbidigo // inform the user about, a real logger is not available at this point.
		fmt.Printf("{\"level\":\"WARN\",\"msg\":\"loading configs from .env: %+v\"}\n", dotEnvErr)
	}
	// Parse environment variabless
	errEnv := env.Parse(&cfg)
	if errEnv != nil {
		//nolint:forbidigo // inform the user about, a real logger is not available at this point.
		fmt.Printf("{\"level\":\"ERROR\",\"msg\":\"error parsing configs from environment: %+v\"}\n", errEnv)

		return nil, errEnv
	}

	return &cfg, nil
}
