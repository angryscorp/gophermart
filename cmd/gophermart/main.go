package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/gophermart/internal/config"
	"github.com/angryscorp/gophermart/internal/http/handler"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/angryscorp/gophermart/internal/repository/migration"
	"github.com/angryscorp/gophermart/internal/repository/users"
	"github.com/angryscorp/gophermart/internal/usecase/auth"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	zeroLogger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	if err := migration.Migrate(cfg.DatabaseDSN); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	usersRepository, err := users.New(cfg.DatabaseDSN)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	authUsecase := auth.New(usersRepository)

	r := router.New(zeroLogger)
	r.RegisterAuth(handler.NewAuth(authUsecase))

	err = r.Run(cfg.ServerAddress)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
