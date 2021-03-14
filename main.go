package main

import (
	"github.com/mp-hl-2021/lenkeforkortelse/api"
	"github.com/mp-hl-2021/lenkeforkortelse/usecases"
	"net/http"
	"time"
)

func main() {
	service := api.NewApi(&usecases.AccountUseCases{}, &usecases.LinkUseCases{})

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
