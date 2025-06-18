package main

import (
	"flag"
	"fmt"
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

	balanceRepository, err := repositoryBalance.New(cfg.DatabaseDSN)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	requestChan := make(chan string)
	responseChan := make(chan model.Accrual)

	accrualAdapter := accrual.NewAdapter(&http.Client{}, zeroLogger, cfg.AccrualAddress)
	accrualWorker := accrual.NewWorker(accrualAdapter, 10, requestChan, responseChan) // TODO
	accrualWorker.Run()

	authUsecase := auth.New(usersRepository, "secret") // TODO
	ordersUsecase := orders.New(ordersRepository, requestChan, responseChan, zeroLogger)
	balanceUsecase := balance.New(balanceRepository)

	r := router.New(zeroLogger, authUsecase)
	r.RegisterAuth(handlerAuth.New(authUsecase))
	r.RegisterOrders(handlerOrders.New(ordersUsecase))
	r.RegisterBalance(handlerBalance.New(balanceUsecase))

	err = r.Run(cfg.ServerAddress)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
