package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mp-hl-2021/chat/usecases"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AccountUseCasesMock struct {}

func (AccountUseCasesMock) CreateAccount(login, password string) (usecases.Account, error) {
	switch login {
	case "alice":
		return usecases.Account{
			Id: "1",
		}, nil
	default:
		return usecases.Account{}, errors.New("failed to create an account")
	}
}

func (AccountUseCasesMock) GetAccountById(id string) (usecases.Account, error) {
	panic("implement me")
}

func (AccountUseCasesMock) LoginToAccount(login, password string) (string, error) {
	if login == "alice" && password == "123" {
		return "token", nil
	}
	return "", errors.New("invalid login or password")
}

func Test_postSignup(t *testing.T) {
	service := NewApi(&AccountUseCasesMock{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signup")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("failed to create account", func(t *testing.T) {
		m := postSignupRequestModel{
			Login:    "bob",
			Password: "123",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusInternalServerError)
	})

	t.Run("successful account creation", func(t *testing.T) {
		m := postSignupRequestModel{
			Login:    "alice",
			Password: "123",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusCreated)

		location := resp.Header().Get("Location")
		if  location != "/accounts/1" {
			t.Errorf("Server MUST return %s Location header, but %s given", "/accounts/1", location)
		}
	})
}

func Test_postSignin(t *testing.T) {
	service := NewApi(&AccountUseCasesMock{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signin")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("failed login with incorrect login or password", func(t *testing.T) {
		m := postSignupRequestModel{
			Login:    "bob",
			Password: "123",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("successful login with correct password", func(t *testing.T) {
		m := postSignupRequestModel{
			Login:    "alice",
			Password: "123",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusOK)
	})
}

func assertStatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Server MUST return %d (%s) status code, but %d (%s) given",
			expectedCode, http.StatusText(expectedCode), actualCode, http.StatusText(actualCode))
	}
}

func invalidJsonTest(router http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{a:")))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}