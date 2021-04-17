package link

import (
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

func (a *LinkUseCases) DeleteLink(link string, accountId string) error {
	err := a.LinkStorage.DeleteLink(link, accountId)
	return err
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
