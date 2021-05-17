package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/account"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/usecases/link"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Api struct {
	AccountUseCases account.AccountUseCasesInterface
	LinkUseCases    link.LinkUseCasesInterface
}

func NewApi(a account.AccountUseCasesInterface, l link.LinkUseCasesInterface) *Api {
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
	// Alya - krasava.
	router.HandleFunc("/link/{link_id}", a.getPage).Methods(http.MethodGet)

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	// lookup all my links
	router.HandleFunc("/accounts/{id}", a.authenticate(a.getAccount)).Methods(http.MethodGet)

	// create link with account
	router.HandleFunc("/accounts/{id}/", a.authenticate(a.postCreateUserLink)).Methods(http.MethodPost)

	// /accounts/{id}/delete/{link_id}
	router.HandleFunc("/accounts/{id}/delete/{link_id}", a.authenticate(a.getDeleteLink)).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())
	//log.Fatalln(http.ListenAndServe(":9090", nil))
	//
	//router.Use(prom.Measurer())
	//router.Use(a.logger)
	//fmt.Println("test")

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
			account.ErrInvalidLoginString,
			account.ErrInvalidPasswordString,
			account.ErrTooShortString,
			account.ErrTooLongString,
			account.ErrNoCapitalLetters,
			account.ErrNoDigits:

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

	shortLink, err := a.LinkUseCases.CutLink(m.Link, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(shortLink)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// getPage handles request for short link, redirect to user's source web page
func (a *Api) getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId, ok := vars["link_id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l, err := a.LinkUseCases.GetLinkByLinkId(linkId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, l, http.StatusSeeOther)
}

type getAccountResponseModel struct {
	Links []link.Link `json:"links"`
}

// getAccount handles request for user's account information.
func (a *Api) getAccount(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value("account_id").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	accountId, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accountId != aid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	links, err := a.LinkUseCases.GetLinksByAccountId(aid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ret := getAccountResponseModel{Links: make([]link.Link, 0, len(links))}
	for _, l := range links {
		ret.Links = append(ret.Links, link.Link{
			LinkId: l.LinkId,
			Link:   l.Link,
		})
	}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type postAccountLinkRequestModel struct {
	Link string `json:"link"`
}

// postCreateUserLink handles request for creating short link from specific user
func (a *Api) postCreateUserLink(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value("account_id").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	accountId, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accountId != aid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var m postAccountLinkRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink, err := a.LinkUseCases.CutLink(m.Link, &aid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(shortLink)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// getDeleteLink handles link deletion request from user
func (a *Api) getDeleteLink(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value("account_id").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	accountId, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accountId != aid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	linkId, ok := vars["link_id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := a.LinkUseCases.DeleteLink(linkId, aid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
