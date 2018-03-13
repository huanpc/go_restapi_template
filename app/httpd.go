package main

import (
	"apistream/router"
	"net/http"
	"fmt"
	"log"

	"apistream/config"
	"github.com/go-chi/chi"
)

func main() {
	// Read config file
	cfg := config.AppConfig()
	r := chi.NewRouter()

	router.Register(r)
	log.Printf(fmt.Sprintf("%v:%v", cfg.HOST_ADDRESS, cfg.Port))

	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", cfg.HOST_ADDRESS, cfg.Port), r); err != nil {
		panic(err)
	}
}
