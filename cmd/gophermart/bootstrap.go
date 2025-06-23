package main

import (
	"github.com/angryscorp/gophermart/internal/config"
	"github.com/angryscorp/gophermart/internal/domain/model"
	handlerAuth "github.com/angryscorp/gophermart/internal/http/handler/auth"
	handlerBalance "github.com/angryscorp/gophermart/internal/http/handler/balance"
	handlerOrders "github.com/angryscorp/gophermart/internal/http/handler/orders"
	"github.com/angryscorp/gophermart/internal/http/router"
	repositoryBalance "github.com/angryscorp/gophermart/internal/repository/balance"
	"github.com/angryscorp/gophermart/internal/repository/migration"
	repositoryOrders "github.com/angryscorp/gophermart/internal/repository/orders"
	repositoryUsers "github.com/angryscorp/gophermart/internal/repository/users"
	"github.com/angryscorp/gophermart/internal/usecase/accrual"
	"github.com/angryscorp/gophermart/internal/usecase/auth"
	"github.com/angryscorp/gophermart/internal/usecase/balance"
	"github.com/angryscorp/gophermart/internal/usecase/orders"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func bootstrap(cfg *config.Config) (*router.Router, error) {
	zeroLogger := zerolog.New(os.Stdout).
		Level(cfg.LogLevel()).
		With().
		Timestamp().
		Logger()

	if err := migration.Migrate(cfg.DatabaseDSN); err != nil {
		return nil, err
	}

	usersRepository, err := repositoryUsers.New(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	ordersRepository, err := repositoryOrders.New(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	balanceRepository, err := repositoryBalance.New(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	requestChan := make(chan string)
	responseChan := make(chan *model.Accrual)

	accrualAdapter := accrual.NewAdapter(&http.Client{}, zeroLogger, cfg.AccrualAddress)
	accrualWorker := accrual.NewWorker(accrualAdapter, cfg.RateLimiter, requestChan, responseChan)
	accrualWorker.Run()

	authUsecase := auth.New(usersRepository, cfg.JWTSecret)
	ordersUsecase := orders.New(ordersRepository, requestChan, responseChan, zeroLogger)
	balanceUsecase := balance.New(balanceRepository)

	r := router.New(zeroLogger, authUsecase)
	r.RegisterAuth(handlerAuth.New(authUsecase, zeroLogger))
	r.RegisterOrders(handlerOrders.New(ordersUsecase, zeroLogger))
	r.RegisterBalance(handlerBalance.New(balanceUsecase, zeroLogger))

	return &r, nil
}
