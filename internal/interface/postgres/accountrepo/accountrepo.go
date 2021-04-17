package accountrepo

import (
	"database/sql"
	"errors"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/account"
	"strconv"
)

type Postgres struct {
	conn *sql.DB
}

var (
	ErrConversion = errors.New("cant typecast string account id to int")
)

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryCreateAccount = `
	INSERT INTO accounts(
	                     login, password
	) VALUES ($1, $2)
	RETURNING id	
`

func (p *Postgres) CreateAccount(cred account.Credentials) (account.Account, error) {
	a := account.Account{Credentials: cred}
	row := p.conn.QueryRow(queryCreateAccount, cred.Login, cred.Password)
	err := row.Scan(&a.Id)
	if err != nil && err == sql.ErrNoRows {
		return account.Account{}, account.ErrAlreadyExist
	}
	return a, err
}

const queryGetAccountById = `
	select id, login, password from accounts where id = $1
`

func (p *Postgres) GetAccountById(id string) (account.Account, error) {
	a := account.Account{}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return a, ErrConversion
	}

	row := p.conn.QueryRow(queryGetAccountById, intId)

	accountId := -1
	err = row.Scan(&accountId, &a.Login, &a.Password)
	a.Id = strconv.Itoa(accountId)

	if err != nil && err == sql.ErrNoRows {
		return a, account.ErrNotFound
	}
	return a, err
}

const queryGetAccountByLogin = `
	select id, login, password from accounts where login = $1
`

func (p *Postgres) GetAccountByLogin(login string) (account.Account, error) {
	a := account.Account{}
	row := p.conn.QueryRow(queryGetAccountByLogin, login)
	err := row.Scan(&a.Id, &a.Login, &a.Password)
	if err != nil && err == sql.ErrNoRows {
		return a, account.ErrNotFound
	}
	return a, err
}
