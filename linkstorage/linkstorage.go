package linkstorage

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrAccessDenied = errors.New("access denied")
)

type Link struct {
	LinkId    string
	Link      string
	AccountId *string
}

type Interface interface {
	CheckIfLinkExists(linkId string) bool
	StoreLink(link Link) (Link, error)
	DeleteLink(linkId string, accountId string) error
	GetLinkByLinkId(linkId string) (Link, error)
	GetLinksByAccountId(accountId string) ([]Link, error)
}
