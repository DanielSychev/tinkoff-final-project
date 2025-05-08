package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GrpcPort int `yaml:"grpc_port" env:"GRPC_PORT"`
	RestPort int `yaml:"rest_port" env:"REST_PORT"`
}

func NewConfig() (*Config, error) {
	c := new(Config)
	if err := cleanenv.ReadConfig("./internal/config/.env", c); err != nil {
		return nil, err
	}
	return c, nil
}
