package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/lenkeforkortelse/usecases"
	"net/http"
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

// NewRouter creates all endpoints for chat app.
func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	// /links post request to create link <link to source>, returns <short link>
	router.HandleFunc("/links", a.postCreateLink).Methods(http.MethodPost)

	// /{link} get redirect
	router.HandleFunc("/{link_id}", a.getPage).Methods(http.MethodGet)

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	// lookup all my links
	router.HandleFunc("/accounts/{id}", a.getAccount).Methods(http.MethodGet)

	// create link with account
	router.HandleFunc("/accounts/{id}/", a.postCreateUserLink).Methods(http.MethodGet)

	// /accounts/{id}/delete/{link_id}
	router.HandleFunc("/accounts/{id}/delete/{link_id}", a.getDeleteLink).Methods(http.MethodGet)

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
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
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
	w.Write([]byte(token))
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

	w.Write([]byte(shortLink))
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
	// todo: authorize user, list all links for this account
	w.WriteHeader(http.StatusNotImplemented)
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
