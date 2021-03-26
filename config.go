package main

import "fmt"

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
	Port int
	Env  string
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		Port: 3000,
		Env:  "dev",
	}
}

// IsProd ...
func (c Config) IsProd() bool {
	return c.Env == "prod"
}

// # models/users.go
// const userPwPepper = "Nhfsf632dfsg"
// const hmacSecretKey = "secret-key"

// # models/services.go
// db, err := gorm.Open("postgres", connectionInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.LogMode(true)
