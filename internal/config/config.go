package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	EnvVarPrefix      = "YANI_"
	ConfigPathEnvVar  = EnvVarPrefix + "CONFIG_PATH"
	DefaultConfigPath = "./config/config.toml"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

func (e Env) String() string { return string(e) }

func (e Env) IsValid() bool {
	switch e {
	case EnvLocal, EnvDev, EnvProd:
		return true
	default:
		return false
	}
}

type Config struct {
	Env     Env `koanf:"env"`
	Storage struct {
		DSN string `koanf:"dsn"`
	} `koanf:"storage"`
	GRPC struct {
		Port    int           `koanf:"port"`
		Timeout time.Duration `koanf:"timeout"`
	} `koanf:"grpc"`
}

func (c *Config) Validate() error {
	if !c.Env.IsValid() {
		return fmt.Errorf("invalid env: %q (allowed: local, dev, prod)", c.Env)
	}
	if c.Storage.DSN == "" {
		return fmt.Errorf("storage dsn is required")
	}
	if c.GRPC.Port <= 0 {
		return fmt.Errorf("invalid grpc port: %d (allowed: > 0)", c.GRPC.Port)
	}
	if c.GRPC.Timeout <= 0 {
		return fmt.Errorf("invalid grpc timeout: %v (allowed: > 0)", c.GRPC.Timeout)
	}
	return nil
}

func Load() (*Config, error) {
	k := koanf.New(".")

	configPath := getConfigPath()
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	normalizeEnvKey := func(s string) string {
		s = strings.TrimPrefix(s, EnvVarPrefix)
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "_", ".")
	}

	if err := k.Load(env.Provider(EnvVarPrefix, ".", normalizeEnvKey), nil); err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}

	c := &Config{}
	if err := k.Unmarshal("", c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %w", configPath, err)
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		panic(err)
	}
	return c
}

func getConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath != "" {
		return configPath
	}

	configPath = os.Getenv(ConfigPathEnvVar)
	if configPath != "" {
		return configPath
	}

	return DefaultConfigPath
}
