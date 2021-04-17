package account

import (
	"github.com/mp-hl-2021/lenkeforkortelse/accountstorage"
	"github.com/mp-hl-2021/lenkeforkortelse/auth"

	"errors"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

var (
	ErrInvalidLoginString    = errors.New("login string contains invalid character")
	ErrInvalidPasswordString = errors.New("password string contains invalid character")
	ErrTooShortString        = errors.New("too short string")
	ErrTooLongString         = errors.New("too long string")
	ErrNoCapitalLetters      = errors.New("password string does not contain capital letters")
	ErrNoDigits              = errors.New("password string does not contain digits")
)

const (
	minLoginLength    = 4
	maxLoginLength    = 50
	minPasswordLength = 6
	maxPasswordLength = 50
)

type Account struct {
	Id string
}

type AccountUseCasesInterface interface {
	CreateAccount(login, password string) (Account, error)
	GetAccountById(id string) (Account, error)
	LoginToAccount(login, password string) (string, error)
	Authenticate(token string) (string, error)
}

type AccountUseCases struct {
	AccountStorage accountstorage.Interface
	Auth           auth.Interface
}

func (a *AccountUseCases) CreateAccount(login, password string) (Account, error) {
	if err := validateLogin(login); err != nil {
		return Account{}, err
	}
	if err := validatePassword(password); err != nil {
		return Account{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	acc, err := a.AccountStorage.CreateAccount(accountstorage.Credentials{
		Login:    login,
		Password: string(hashedPassword),
	})
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, nil
}

func (a *AccountUseCases) GetAccountById(id string) (Account, error) {
	acc, err := a.AccountStorage.GetAccountById(id)
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, err
}

func (a *AccountUseCases) LoginToAccount(login, password string) (string, error) {
	if err := validateLogin(login); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}
	acc, err := a.AccountStorage.GetAccountByLogin(login)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", err
	}
	token, err := a.Auth.IssueToken(acc.Id)
	if err != nil {
		return "", err
	}
	return token, err
}

func (a *AccountUseCases) Authenticate(token string) (string, error) {
	return a.Auth.UserIdByToken(token)
}

func validateLogin(login string) error {
	chars := 0
	for _, r := range login {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidLoginString
		}
		chars++
	}
	if chars < minLoginLength {
		return ErrTooShortString
	}
	if chars > maxLoginLength {
		return ErrTooLongString
	}
	return nil
}

func validatePassword(password string) error {
	chars := 0
	digits := 0
	capitalLetters := 0
	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return ErrInvalidPasswordString
		}

		if unicode.IsDigit(r) {
			digits += 1
		}
		if unicode.IsUpper(r) {
			capitalLetters += 1
		}

		chars++
	}
	if chars < minPasswordLength {
		return ErrTooShortString
	}
	if chars > maxPasswordLength {
		return ErrTooLongString
	}
	if digits == 0 {
		return ErrNoDigits
	}
	if capitalLetters == 0 {
		return ErrNoCapitalLetters
	}
	return nil
}
