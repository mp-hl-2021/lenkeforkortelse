package main

import (
	"flag"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/httpapi"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/memory/accountrepo"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/memory/linkrepo"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/service/token"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/account"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/link"
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

	a, err := token.NewJwtHandler(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	accountUseCases := &account.AccountUseCases{
		AccountStorage: accountrepo.NewMemory(),
		Auth:           a,
	}

	linkUseCases := &link.LinkUseCases{
		LinkStorage: linkrepo.NewMemory(),
	}

	service := httpapi.NewApi(accountUseCases, linkUseCases)

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
