package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/httpapi"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/postgres/accountrepo"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/interface/postgres/linkrepo"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/pipeline"
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

	// todo: pass connection args through config
	connStr := "user=postgres password=12345678 host=db dbname=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//defer conn.Close()

	accountUseCases := &account.AccountUseCases{
		AccountStorage: accountrepo.New(conn),
		Auth:           a,
	}

	linkUseCases := &link.LinkUseCases{
		LinkStorage: linkrepo.New(conn),
	}

	pipeline.LinkStatusUpdater(linkUseCases)

	service := httpapi.NewApi(accountUseCases, linkUseCases)

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
