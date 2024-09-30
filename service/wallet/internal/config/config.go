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
	Tracer            trace.Config
	ServiceName       string `env:"SERVICE_NAME,default=wallet-server"`
	AppEnv            string `env:"APP_ENV,default=development"`
	Port              string `env:"PORT,default=8004"`
	PrometheusPort    string `env:"PROMETHEUS_PORT,default=7004"`
	Username          string `env:"USERNAME,default=wallet-user"`
	Password          string `env:"PASSWORD,default=wallet-password"`
	AppliedAuthBearer string `env:"APPLIED_AUTH_BEARER"`
	AppliedAuthBasic  string `env:"APPLIED_AUTH_BASIC"`
	SecretKey         string `env:"TOKEN_SECRET_KEY,required"`
	Redis             Redis
	Postgres          sdkpg.Config
}

// Redis holds configuration for Redis.
type Redis struct {
	Address string `env:"REDIS_ADDRESS,default=localhost:6379"`
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
