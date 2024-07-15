package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
)

// Config holds configuration for the project.
type Config struct {
	Postgres       sdkpg.Config
	Tracer         trace.Config
	ServiceName    string `env:"SERVICE_NAME,default=auth-server"`
	AppEnv         string `env:"APP_ENV,default=development"`
	Port           string `env:"PORT,default=8002"`
	PrometheusPort string `env:"PROMETHEUS_PORT,default=7002"`
	Keycloak       Keycloak
}

// Keycloak holds configuration for Keycloak.
type Keycloak struct {
	Address string `env:"KEYCLOAK_ADDRESS,default=http://localhost:8080/"`
	Realm   string `env:"KEYCLOAK_REALM,required"`
	Timeout int    `env:"KEYCLOAK_TIMEOUT_SECONDS,default=5"`
}

// NewConfig creates an instance of Config.
// It needs the path of the env file to be used.
func NewConfig(env string) (*Config, error) {
	// just skip loading env files if it is not exists, env files only used in local dev
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "[NewConfig] error decoding env")
	}

	return &config, nil
}
