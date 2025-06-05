package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN    string `env:"DATABASE_URI"`
}

func NewConfig() (Config, error) {
	serverAddress := flag.String("a", "localhost:8081", "Gophermart server address (default: localhost:8081)")
	accrualAddress := flag.String("r", "localhost:8080", "Accrual System address (default: localhost:8080)")
	databaseDSN := flag.String("d", "", "Database DSN (default: empty, mandatory)")

	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return Config{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	config := Config{
		ServerAddress:  *serverAddress,
		AccrualAddress: *accrualAddress,
		DatabaseDSN:    *databaseDSN,
	}

	// ENV vars
	err := env.Parse(&config)
	if err != nil {
		return Config{}, err
	}

	// Mandatory fields
	if config.DatabaseDSN == "" {
		return Config{}, fmt.Errorf("database DSN is mandatory")
	}

	return config, nil
}
