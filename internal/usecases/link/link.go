package link

import (
	"fmt"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/link"
	"math/rand"
	"time"
	"unsafe"
)

const (
	linkLength    = 6
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Link struct {
	LinkId string
	Link   string
}

type LinkUseCases struct {
	LinkStorage link.Interface
}

type LinkUseCasesInterface interface {
	GetLinkByLinkId(linkId string) (string, error)
	CutLink(link string, accountId *string) (string, error)
	DeleteLink(linkId string, accountId string) error
	GetLinksByAccountId(accountId string) ([]Link, error)

	//Logging
	LoggerGetLinkByLinkId(
		getLinkByLinkId func(linkId string) (string, error)) func(linkId string) (string, error)
	LoggerCutLink(
		cutLink func(link string, accountId *string) (string, error)) func(link string, accountId *string) (string, error)
	LoggerDeleteLink(
		deleteLink func(linkId string, accountId string) error) func(linkId string, accountId string) error
	LoggerGetLinksByAccountId(
		getLinksByAccountId func(accountId string) ([]Link, error)) func(accountId string) ([]Link, error)
}

func (a *LinkUseCases) GetLinkByLinkId(lnk string) (string, error) {
	l, err := a.LinkStorage.GetLinkByLinkId(lnk)
	if err != nil {
		return "", err
	}
	return l.Link, nil
}

func (a *LinkUseCases) CutLink(lnk string, accountId *string) (string, error) {
	linkId := a.generateFreeLinkId()
	l, err := a.LinkStorage.StoreLink(link.Link{
		LinkId:    linkId,
		Link:      lnk,
		AccountId: accountId,
	})
	if err != nil {
		return "", err
	}
	return l.LinkId, nil
}

func (a *LinkUseCases) DeleteLink(lnk string, accountId string) error {
	dbLink, err := a.LinkStorage.GetLinkByLinkId(lnk)
	if err != nil {
		if err == link.ErrNotFound {
			return nil
		}
		return err
	}
	if dbLink.AccountId != nil && *dbLink.AccountId != accountId {
		return link.ErrAccessDenied
	}
	return a.LinkStorage.DeleteLink(lnk)
}

func (a *LinkUseCases) GetLinksByAccountId(accountId string) ([]Link, error) {
	links, err := a.LinkStorage.GetLinksByAccountId(accountId)
	if err != nil {
		return nil, err
	}
	res := make([]Link, 0, len(links))
	for _, l := range links {
		res = append(res, Link{
			LinkId: l.LinkId,
			Link:   l.Link,
		})
	}
	return res, nil
}

var src = rand.NewSource(time.Now().UnixNano())

func generateLinkId() (linkId string) {

	b := make([]byte, linkLength)
	for i, cache, remain := linkLength-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (a *LinkUseCases) generateFreeLinkId() (linkId string) {
	for {
		linkId = generateLinkId()
		if !a.LinkStorage.CheckIfLinkExists(linkId) {
			break
		}
	}
	return linkId
}

func (a *LinkUseCases) logger(method string, err error, start time.Time) {
	status := "SUCCESS"
	if err != nil {
		status = err.Error()
	}
	fmt.Printf("method: %s; status-code: %s; call time: %v; duration: %v;\n",
		method, status, start, time.Since(start))
}

func (a *LinkUseCases) LoggerGetLinkByLinkId(
	getLinkByLinkId func(linkId string) (string, error)) func(linkId string) (string, error) {

	return func(linkId string) (string, error) {
		start := time.Now()
		link, err := getLinkByLinkId(linkId)
		a.logger("GetLinkByLinkId", err, start)
		return link, err
	}
}

func (a *LinkUseCases) LoggerCutLink(
	cutLink func(link string, accountId *string) (string, error)) func(link string, accountId *string) (string, error) {

	return func(link string, accountId *string) (string, error) {
		start := time.Now()
		linkId, err := cutLink(link, accountId)
		a.logger("CutLink", err, start)
		return linkId, err
	}
}

func (a *LinkUseCases) LoggerDeleteLink(
	deleteLink func(linkId string, accountId string) error) func(linkId string, accountId string) error {

	return func(linkId string, accountId string) error {
		start := time.Now()
		err := deleteLink(linkId, accountId)
		a.logger("DeleteLink", err, start)
		return err
	}
}

func (a *LinkUseCases) LoggerGetLinksByAccountId(
	getLinksByAccountId func(accountId string) ([]Link, error)) func(accountId string) ([]Link, error) {

	return func(accountId string) ([]Link, error) {
		start := time.Now()
		res, err := getLinksByAccountId(accountId)
		a.logger("GetLinksByAccountId", err, start)
		return res, err
	}
}