package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	Env string `yaml:"env"` // Environment: dev or prod
}

func expandEnvVars(content string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)
	return re.ReplaceAllStringFunc(content, func(s string) string {
		key := re.FindStringSubmatch(s)[1]
		return os.Getenv(key)
	})
}

func LoadConfig() *Config {
	cfg := &Config{}
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to read config.yaml: %v", err))
	}
	// Expand environment variables before unmarshaling
	expanded := expandEnvVars(string(file))
	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		panic(fmt.Sprintf("Failed to parse config.yaml: %v", err))
	}
	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=prefer connect_timeout=5",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name,
	)
}
