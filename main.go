package main

import (
	"github.com/mp-hl-2021/chat/api"
	"github.com/mp-hl-2021/chat/usecases"
	"net/http"
	"time"
)

func main() {
	accountUseCases := &usecases.AccountUseCases{}

	service := api.NewApi(accountUseCases)

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

