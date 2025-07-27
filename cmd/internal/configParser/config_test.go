package configparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCfgValue_String_Exists(t *testing.T) {
	os.Setenv("MY_KEY", "hello")
	val := GetCfgValue("MY_KEY", "default")
	assert.Equal(t, "hello", val)
}

func TestGetCfgValue_String_NotExists(t *testing.T) {
	os.Unsetenv("MY_KEY")
	val := GetCfgValue("MY_KEY", "default")
	assert.Equal(t, "default", val)
}

func TestGetCfgValue_Int_Valid(t *testing.T) {
	os.Setenv("INT_KEY", "42")
	val := GetCfgValue("INT_KEY", 10)
	assert.Equal(t, 42, val)
}

func TestGetCfgValue_Int_Invalid(t *testing.T) {
	os.Setenv("INT_KEY", "notanint")
	val := GetCfgValue("INT_KEY", 10)
	assert.Equal(t, 10, val)
}

func TestGetCfgValue_Int_NotExists(t *testing.T) {
	os.Unsetenv("INT_KEY")
	val := GetCfgValue("INT_KEY", 99)
	assert.Equal(t, 99, val)
}

func TestLoadConfig_Valid(t *testing.T) {
	yaml := `
db:
  host: localhost
  port: "5432"
  name: testdb
  user: user
  password: pass
  readiness:
    max_wait_seconds: 10
    poll_interval_seconds: 2

kafka:
  broker: kafka:9092
  readiness:
    max_wait_seconds: 5
    poll_interval_seconds: 1

http:
  addr: 0.0.0.0
  port: "8080"
  read_timeout_seconds: 30
  write_timeout_seconds: 30
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yaml)
	assert.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, 10, cfg.DB.Readiness.MaxWaitSeconds)
	assert.Equal(t, "kafka:9092", cfg.Kafka.Broker)
	assert.Equal(t, "0.0.0.0", cfg.HTTP.Addr)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	cfg, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
	assert.NotNil(t, cfg)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(":::bad_yaml:::")
	assert.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	assert.Error(t, err)
	assert.NotNil(t, cfg)
}
