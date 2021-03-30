package api

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/lenkeforkortelse/usecases"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Api struct {
	AccountUseCases usecases.AccountUseCasesInterface
	LinkUseCases    usecases.LinkUseCasesInterface
}

func NewApi(a usecases.AccountUseCasesInterface, l usecases.LinkUseCasesInterface) *Api {
	return &Api{
		AccountUseCases: a,
		LinkUseCases:    l,
	}
}

// NewRouter creates all endpoints for app.
func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	// /links post request to create link <link to source>, returns <short link>
	router.HandleFunc("/links", a.postCreateLink).Methods(http.MethodPost)

	// /{link} get redirect
	router.HandleFunc("/{link_id}", a.getPage).Methods(http.MethodGet)

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	// lookup all my links
	router.HandleFunc("/accounts/{id}", a.authorize(a.getAccount)).Methods(http.MethodGet)

	// create link with account
	router.HandleFunc("/accounts/{id}/", a.authorize(a.postCreateUserLink)).Methods(http.MethodGet)

	// /accounts/{id}/delete/{link_id}
	router.HandleFunc("/accounts/{id}/delete/{link_id}", a.authorize(a.getDeleteLink)).Methods(http.MethodGet)

	return router
}

type postSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// postSignup handles request for a new account creation.
func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := a.AccountUseCases.CreateAccount(m.Login, m.Password)
	if err != nil {
		switch err {
		case
			usecases.ErrInvalidLoginString,
			usecases.ErrInvalidPasswordString,
			usecases.ErrTooShortString,
			usecases.ErrTooLongString,
			usecases.ErrNoCapitalLetters,
			usecases.ErrNoDigits:

			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println(err)
		return
	}

	location := fmt.Sprintf("/accounts/%s", acc.Id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

// postSignin handles login request for existing user.
func (a *Api) postSignin(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(m.Login, m.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	if _, err := w.Write([]byte(token)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type postLinkRequestModel struct {
	Link string `json:"link"`
}

// postCreateLink handles creating short link from user's link
func (a *Api) postCreateLink(w http.ResponseWriter, r *http.Request) {
	var m postLinkRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink, err := a.LinkUseCases.CutLink(m.Link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(shortLink)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type getPageResponseModel struct {
	Link string `json:"link"`
}

// getPage handles request for short link, redirect to user's source web page
func (a *Api) getPage(w http.ResponseWriter, r *http.Request) {
	// todo: use mux.Var(r) to parse {link_id} into m.Link
	w.WriteHeader(http.StatusNotImplemented)
}

type getAccountResponseModel struct {
	Id string `json:"id"`
}

// getAccount handles request for user's account information.
func (a *Api) getAccount(w http.ResponseWriter, r *http.Request) {
	// todo: list all links for this account

	// Following implementation was added just to check that the authentication works correctly and the URLs are protected
	if _, err := w.Write([]byte("Hi!")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type postAccountLinkRequestModel struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

// postCreateUserLink handles request for creating short link from specific user
func (a *Api) postCreateUserLink(w http.ResponseWriter, r *http.Request) {
	// todo: authorize user
	// todo: use mux.Var(r) to parse {link_id} into m.Link, and {id} into m.Id
	w.WriteHeader(http.StatusNotImplemented)
}

type getLinkDeleteRequestModel struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

// getDeleteLink handles link deletion request from user
func (a *Api) getDeleteLink(w http.ResponseWriter, r *http.Request) {
	// todo: authorize user
	// todo: use mux.Var(r) to parse {link_id} into m.Link, and {id} into m.Id
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearHeader := r.Header.Get("Authorization")
		strArr := strings.Split(bearHeader, " ")
		if len(strArr) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := strArr[1]
		id, err := a.AccountUseCases.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "account_id", id)
		handler(w, r.WithContext(ctx))
	}
}
