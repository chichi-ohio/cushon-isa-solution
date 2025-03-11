package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Queue    QueueConfig    `mapstructure:"queue"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type QueueConfig struct {
	Type          string   `mapstructure:"type"`
	BufferSize    int      `mapstructure:"buffer_size"`
	NumWorkers    int      `mapstructure:"num_workers"`
	Brokers       []string `mapstructure:"brokers"`
	Topic         string   `mapstructure:"topic"`
	ConsumerGroup string   `mapstructure:"consumer_group"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set default values
	if config.Queue.BufferSize == 0 {
		config.Queue.BufferSize = 100
	}
	if config.Queue.NumWorkers == 0 {
		config.Queue.NumWorkers = 2
	}
	if len(config.Queue.Brokers) == 0 {
		config.Queue.Brokers = []string{"localhost:9092"}
	}
	if config.Queue.Topic == "" {
		config.Queue.Topic = "investments"
	}
	if config.Queue.ConsumerGroup == "" {
		config.Queue.ConsumerGroup = "investment-processor"
	}
	if config.Logger.Level == "" {
		config.Logger.Level = "info"
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8081)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")

	// Queue defaults
	viper.SetDefault("queue.type", "memory")
	viper.SetDefault("queue.buffer_size", 100)
	viper.SetDefault("queue.num_workers", 2)
	viper.SetDefault("queue.brokers", []string{"localhost:9092"})
	viper.SetDefault("queue.topic", "investments")
	viper.SetDefault("queue.consumer_group", "investment-processor")

	// Logger defaults
	viper.SetDefault("logger.level", "info")
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
