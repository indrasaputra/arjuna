package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Config holds configuration for the project.
type Config struct {
	ServiceName        string `env:"SERVICE_NAME,default=grpc-gateway"`
	AppEnv             string `env:"APP_ENV,default=development"`
	Port               string `env:"PORT,default=8000"`
	UserServiceAddress string `env:"USER_SERVICE_ADDRESS,required"`
	AuthServiceAddress string `env:"AUTH_SERVICE_ADDRESS,required"`
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
