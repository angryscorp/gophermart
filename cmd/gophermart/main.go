package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/angryscorp/gophermart/internal/config"
	_ "github.com/rs/zerolog"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		configJSON, _ := json.Marshal(cfg)
		_, _ = w.Write([]byte(string(configJSON)))
	})

	_ = http.ListenAndServe(":8081", nil)
}
