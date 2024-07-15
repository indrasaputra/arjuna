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
	Postgres                    sdkpg.Config
	Tracer                      trace.Config
	Temporal                    Temporal
	ServiceName                 string `env:"SERVICE_NAME,default=user-server"`
	AppEnv                      string `env:"APP_ENV,default=development"`
	Port                        string `env:"PORT,default=8001"`
	PrometheusPort              string `env:"PROMETHEUS_PORT,default=7001"`
	AuthServiceHost             string `env:"AUTH_SERVICE_HOST,required"`
	Keycloak                    Keycloak
	RelayerSleepTimeMillisecond int `env:"RELAYER_SLEEP_TIME_MILLISECONDS,default=1000"`
}

// Keycloak holds configuration for Keycloak.
type Keycloak struct {
	Address       string `env:"KEYCLOAK_ADDRESS,default=http://localhost:8080/"`
	Realm         string `env:"KEYCLOAK_REALM,required"`
	AdminUser     string `env:"KEYCLOAK_ADMIN_USER,required"`
	AdminPassword string `env:"KEYCLOAK_ADMIN_PASSWORD,required"`
	Timeout       int    `env:"KEYCLOAK_TIMEOUT_SECONDS,default=5"`
}

// Temporal holds configuration for Temporal.
type Temporal struct {
	Address string `env:"TEMPORAL_ADDRESS,default=localhost:7233"`
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
