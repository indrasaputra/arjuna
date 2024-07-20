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
	SecretKey                   string `env:"TOKEN_SECRET_KEY,required"`
	Username                    string `env:"USERNAME,default=user-user"`
	ServiceName                 string `env:"SERVICE_NAME,default=user-server"`
	AppEnv                      string `env:"APP_ENV,default=development"`
	Port                        string `env:"PORT,default=8001"`
	PrometheusPort              string `env:"PROMETHEUS_PORT,default=7001"`
	AuthServiceHost             string `env:"AUTH_SERVICE_HOST,required"`
	WalletServiceHost           string `env:"WALLET_SERVICE_HOST,required"`
	Temporal                    Temporal
	Redis                       Redis
	Password                    string `env:"PASSWORD,default=user-password"`
	AppliedAuthBearer           string `env:"APPLIED_AUTH_BEARER"`
	AppliedAuthBasic            string `env:"APPLIED_AUTH_BASIC"`
	WalletServicePassword       string `env:"WALLET_SERVICE_PASSWORD"`
	AuthServiceUsername         string `env:"AUTH_SERVICE_USERNAME"`
	AuthServicePassword         string `env:"AUTH_SERVICE_PASSWORD"`
	WalletServiceUsername       string `env:"WALLET_SERVICE_USERNAME"`
	RelayerSleepTimeMillisecond int    `env:"RELAYER_SLEEP_TIME_MILLISECONDS,default=1000"`
}

// Temporal holds configuration for Temporal.
type Temporal struct {
	Address string `env:"TEMPORAL_ADDRESS,default=localhost:7233"`
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
