package main

import (
	"github.com/mp-hl-2021/lenkeforkortelse/accountstorage"
	"github.com/mp-hl-2021/lenkeforkortelse/api"
	"github.com/mp-hl-2021/lenkeforkortelse/auth"
	"github.com/mp-hl-2021/lenkeforkortelse/linkstorage"
	"github.com/mp-hl-2021/lenkeforkortelse/usecases/account"
	"github.com/mp-hl-2021/lenkeforkortelse/usecases/link"

	"flag"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := auth.NewJwtHandler(privateKeyBytes, publicKeyBytes, 100 * time.Minute)
	if err != nil {
		panic(err)
	}

	accountUseCases := &account.AccountUseCases{
		AccountStorage: accountstorage.NewMemory(),
		Auth: a,
	}

	linkUseCases := &link.LinkUseCases{
		LinkStorage: linkstorage.NewMemory(),
	}

	service := api.NewApi(accountUseCases, linkUseCases)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
