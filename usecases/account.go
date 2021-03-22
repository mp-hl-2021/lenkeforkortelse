package usecases

type Account struct {
	Id string
}

type AccountUseCasesInterface interface {
	CreateAccount(login, password string) (Account, error)
	GetAccountById(id string) (Account, error)
	LoginToAccount(login, password string) (string, error)
}

type AccountUseCases struct{}

func (AccountUseCases) CreateAccount(login, password string) (Account, error) {
	panic("implement me")
}

func (AccountUseCases) GetAccountById(id string) (Account, error) {
	panic("implement me")
}

func (AccountUseCases) LoginToAccount(login, password string) (string, error) {
	panic("implement me")
}



