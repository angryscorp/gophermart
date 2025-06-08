package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/gophermart/internal/config"
	handlerAuth "github.com/angryscorp/gophermart/internal/http/handler/auth"
	handlerOrders "github.com/angryscorp/gophermart/internal/http/handler/orders"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/angryscorp/gophermart/internal/repository/migration"
	repositoryOrders "github.com/angryscorp/gophermart/internal/repository/orders"
	repositoryUsers "github.com/angryscorp/gophermart/internal/repository/users"
	"github.com/angryscorp/gophermart/internal/usecase/auth"
	"github.com/angryscorp/gophermart/internal/usecase/orders"
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

	usersRepository, err := repositoryUsers.New(cfg.DatabaseDSN)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	ordersRepository, err := repositoryOrders.New(cfg.DatabaseDSN)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	authUsecase := auth.New(usersRepository)
	ordersUsecase := orders.New(ordersRepository)

	r := router.New(zeroLogger)
	r.RegisterAuth(handlerAuth.New(authUsecase))
	r.RegisterOrders(handlerOrders.New(ordersUsecase))

	err = r.Run(cfg.ServerAddress)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
