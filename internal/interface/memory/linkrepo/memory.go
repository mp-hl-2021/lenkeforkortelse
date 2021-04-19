package linkrepo

import (
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/link"
	"sync"
)

type Memory struct {
	linkByLinkId     map[string]link.Link
	linksByAccountId map[string]map[string]link.Link
	mu               *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		linkByLinkId:     make(map[string]link.Link),
		linksByAccountId: make(map[string]map[string]link.Link),
		mu:               &sync.Mutex{},
	}
}

func (m *Memory) CheckIfLinkExists(linkId string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.linkByLinkId[linkId]
	return ok
}

func (m *Memory) StoreLink(lnk link.Link) (link.Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.linkByLinkId[lnk.LinkId]; ok {
		return link.Link{}, link.ErrAlreadyExist
	}
	m.linkByLinkId[lnk.LinkId] = lnk
	if lnk.AccountId != nil {
		links, ok := m.linksByAccountId[*lnk.AccountId]
		if !ok {
			links = make(map[string]link.Link)
		}
		links[lnk.LinkId] = lnk
		m.linksByAccountId[*lnk.AccountId] = links
	}
	return lnk, nil
}

func (m *Memory) DeleteLink(lnk string, accountId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.linkByLinkId[lnk]
	if !ok {
		return link.ErrNotFound
	}
	if l.AccountId == nil || *l.AccountId != accountId {
		return link.ErrAccessDenied
	}
	delete(m.linkByLinkId, lnk)
	delete(m.linksByAccountId[*l.AccountId], lnk)
	return nil
}

func (m *Memory) GetLinkByLinkId(lnk string) (link.Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.linkByLinkId[lnk]
	if !ok {
		return link.Link{}, link.ErrNotFound
	}
	return l, nil
}

func (m *Memory) GetLinksByAccountId(accountId string) ([]link.Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lnks, ok := m.linksByAccountId[accountId]
	if !ok {
		lnks = make(map[string]link.Link)
		m.linksByAccountId[accountId] = lnks
	}
	links := make([]link.Link, 0, len(lnks))
	for _, val := range lnks {
		links = append(links, val)
	}
	return links, nil
}
