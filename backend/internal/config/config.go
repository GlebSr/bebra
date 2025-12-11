package config

import (
	"context"
	"encoding/json"
	"os"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"gopkg.in/yaml.v2"
)

// JWTConfig содержит настройки для JWT.
type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
	TokenTTL  int    `yaml:"token_ttl"`
}

// RefreshTokenConfig содержит настройки для токенов обновления.
type RefreshTokenConfig struct {
	SecretKey string `yaml:"secret_key"`
	TokenTTL  int    `yaml:"token_ttl"`
}

// DatabaseConfig содержит настройки для подключения к базе данных.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// GetDBUrl формирует строку подключения к базе данных PostgreSQL.
func (d *DatabaseConfig) GetDBUrl() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.DBName + "?sslmode=disable"
}

// ServerConfig содержит настройки для сервера.
type ServerConfig struct {
	Port string `yaml:"port"`
}

// AppConfig определяет интерфейс для работы с конфигурацией приложения.
type AppConfig interface {
	GetJsonConfig(ctx context.Context) (string, error)
	UpdateJsonConfig(ctx context.Context, config string) error
	GetDatabaseConfig() DatabaseConfig
	GetServerConfig() ServerConfig
	GetJWTConfig() JWTConfig
	GetRefreshTokenConfig() RefreshTokenConfig
}

// Config представляет собой основную структуру конфигурации приложения.
type Config struct {
	Database DatabaseConfig     `yaml:"database"`
	Server   ServerConfig       `yaml:"server"`
	JWT      JWTConfig          `yaml:"jwt"`
	Refresh  RefreshTokenConfig `yaml:"refresh"`
}

// LoadConfig загружает конфигурацию из YAML-файла.
func LoadConfig(ctx context.Context, configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		logger.Errorf(ctx, "Error opening config file: %v", err)

		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		logger.Errorf(ctx, "Error decoding config file: %v", err)

		return nil, err
	}

	return config, nil
}

// GetJsonConfig возвращает конфигурацию в формате JSON.
func (c *Config) GetJsonConfig(ctx context.Context) (string, error) {
	jsonConfig, err := json.Marshal(c)
	if err != nil {
		logger.Errorf(ctx, "Error marshalling config to JSON: %v", err)

		return "", err
	}

	return string(jsonConfig), nil
}

// UpdateJsonConfig обновляет конфигурацию из строки JSON.
func (c *Config) UpdateJsonConfig(ctx context.Context, config string) error {
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		logger.Errorf(ctx, "Error unmarshalling JSON config: %v", err)

		return err
	}

	return nil
}

// GetDatabaseConfig возвращает настройки базы данных.
func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

// GetServerConfig возвращает настройки сервера.
func (c *Config) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *Config) GetJWTConfig() JWTConfig {
	return c.JWT
}

func (c *Config) GetRefreshTokenConfig() RefreshTokenConfig {
	return c.Refresh
}
