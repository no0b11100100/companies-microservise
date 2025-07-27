package configparser

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Rediness struct {
	MaxWaitSeconds      int `yaml:"max_wait_seconds"`
	PollIntervalSeconds int `yaml:"poll_interval_seconds"`
}

type DB struct {
	Host      string   `yaml:"host"`
	Port      string   `yaml:"port"`
	Name      string   `yaml:"name"`
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
	Readiness Rediness `yaml:"readiness"`
}

type Kafka struct {
	Broker    string   `yaml:"broker"`
	Readiness Rediness `yaml:"readiness"`
}

type HTTP struct {
	Addr                string `yaml:"addr"`
	Port                string `yaml:"port"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
}

type Config struct {
	DB    DB    `yaml:"db"`
	Kafka Kafka `yaml:"kafka"`
	HTTP  HTTP  `yaml:"http"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return &Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}

func GetCfgValue[T int | string](key string, fallbackValue T) T {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallbackValue
	}

	var value T

	switch any(value).(type) {
	case int:
		parsed, err := strconv.Atoi(valStr)
		if err != nil {
			return fallbackValue
		}
		return any(parsed).(T)
	case string:
		return any(valStr).(T)
	default:
		return fallbackValue
	}
}
