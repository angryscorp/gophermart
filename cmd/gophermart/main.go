package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/gophermart/internal/config"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	server, err := bootstrap(&cfg)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = server.Run(cfg.ServerAddress)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
