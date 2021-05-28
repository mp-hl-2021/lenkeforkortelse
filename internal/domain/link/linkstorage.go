package link

import (
	"errors"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/status"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrAccessDenied = errors.New("access denied")
)

type Link struct {
	LinkId     string
	Link       string
	LinkStatus status.LinkStatus
	AccountId  *string
}

type Interface interface {
	CheckIfLinkExists(linkId string) bool
	StoreLink(link Link) (Link, error)
	DeleteLink(linkId string) error
	GetLinkByLinkId(linkId string) (Link, error)
	GetLinksByAccountId(accountId string) ([]Link, error)
	UpdateLinkStatusByLinkId(linkId string, linkStatus status.LinkStatus) error
	GetAllUserLinks() ([]Link, error)
}
