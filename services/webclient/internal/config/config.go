package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort        int    `env:"ODJ_EE_HTTP_PORT" envDefault:"8080"`
	TargetURL       string `env:"ODJ_DEP_TARGET_URL" envDefault:"http://localhost:8080/api/load"`
	RequestInterval int    `env:"ODJ_DEP_REQUEST_INTERVAL" envDefault:"0"`
	RequestCount    int    `env:"ODJ_DEP_REQUEST_COUNT" envDefault:"0"`
	Mode            string `env:"ODJ_DEP_MODE" envDefault:"unknown"`
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
