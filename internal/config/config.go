package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"homework9/internal/adapters/adrepo/postgres"
)

type Config struct {
	GrpcPort int               `env:"GRPC_PORT" envDefault:"8080"`
	RestPort int               `env:"REST_PORT" envDefault:"8081"`
	PgConfig postgres.PgConfig `env:"POSTGRES"`
}

func NewConfig() (*Config, error) {
	c := new(Config)
	if err := cleanenv.ReadConfig("./internal/config/.env", c); err != nil {
		return nil, err
	}
	return c, nil
}
