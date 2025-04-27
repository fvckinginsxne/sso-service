package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env   string     `env:"APP_ENV" env-required:"true"`
	DB    DBConfig   `env-prefix:"DB_" env-required:"true"`
	GRPC  GRPCConfig `env-prefix:"GRPC_" env-required:"true"`
	Token Token      `env-prefix:"TOKEN_" env-required:"true"`
}

type DBConfig struct {
	Host       string `env:"HOST" env-default:"localhost"`
	Port       string `env:"PORT" env-default:"5432"`
	DockerPort string `env:"DOCKER_PORT" env-default:"5432"`
	User       string `env:"USER" env-required:"true"`
	Password   string `env:"PASSWORD" env-required:"true"`
	Name       string `env:"NAME" env-required:"true"`
}

type GRPCConfig struct {
	Port       int           `env:"PORT" env-default:"50051"`
	DockerPort int           `env:"DOCKER_PORT" env-default:"50051"`
	Timeout    time.Duration `env:"TIMEOUT" env-default:"60s"`
}

type Token struct {
	TTL    time.Duration `env:"TTL" env-required:"true"`
	Secret string        `env:"SECRET" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); err != nil {
		panic("config file does not exist: " + err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
