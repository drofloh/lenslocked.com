package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// PostgresConfig ...
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// DefaultPostgresConfig ...
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "andrewholford",
		Password: "blablabla",
		Name:     "lenslocked_dev",
	}
}

// Dialect ...
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionInfo ...
func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)

}

// Config ...
type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "Nhfsf632dfsg",
		HMACKey:  "secret-key",
		Database: DefaultPostgresConfig(),
	}
}

type MailgunConfig struct {
	APIKey       string `json:"api_key,omitempty"`
	PublicAPIKey string `json:"public_api_key,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

// IsProd ...
func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(configReq bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .config")
	return c
}
