package config

import (
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Clerk  ClerkConfig  `mapstructure:"clerk"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// ClerkConfig holds Clerk authentication configuration
type ClerkConfig struct {
	JWKSURL   string `mapstructure:"jwks_url"`
	SecretKey string `mapstructure:"secret_key"`
}

// Load loads the configuration from file and environment
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/dental-prosthesis/")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("clerk.jwks_url", "")
	viper.SetDefault("clerk.secret_key", "")

	// Environment variables
	viper.SetEnvPrefix("DENTAL")
	viper.AutomaticEnv()

	// Bind specific environment variables
	_ = viper.BindEnv("server.port", "DENTAL_SERVER_PORT")
	_ = viper.BindEnv("server.host", "DENTAL_SERVER_HOST")
	_ = viper.BindEnv("clerk.jwks_url", "CLERK_JWKS_URL")
	_ = viper.BindEnv("clerk.secret_key", "CLERK_SECRET_KEY")

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found is OK - use defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

