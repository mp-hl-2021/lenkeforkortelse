package linkstorage

import (
	"sync"
)

type Memory struct {
	linkByLinkId     map[string]Link
	linksByAccountId map[string]map[string]Link
	mu               *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		linkByLinkId:     make(map[string]Link),
		linksByAccountId: make(map[string]map[string]Link),
		mu:               &sync.Mutex{},
	}
}

func (m *Memory) CheckIfLinkExists(linkId string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.linkByLinkId[linkId]
	return ok
}

func (m *Memory) StoreLink(lnk Link) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.linkByLinkId[lnk.LinkId]; ok {
		return Link{}, ErrAlreadyExist
	}
	m.linkByLinkId[lnk.LinkId] = lnk
	if lnk.AccountId != nil {
		links, ok := m.linksByAccountId[*lnk.AccountId]
		if !ok {
			links = make(map[string]Link)
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
		return ErrNotFound
	}
	if l.AccountId == nil || *l.AccountId != accountId {
		return ErrAccessDenied
	}
	delete(m.linkByLinkId, lnk)
	delete(m.linksByAccountId[*l.AccountId], lnk)
	return nil
}

func (m *Memory) GetLinkByLinkId(lnk string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.linkByLinkId[lnk]
	if !ok {
		return Link{}, ErrNotFound
	}
	return l, nil
}

func (m *Memory) GetLinksByAccountId(accountId string) ([]Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	lnks, ok := m.linksByAccountId[accountId]
	if !ok {
		lnks = make(map[string]Link)
		m.linksByAccountId[accountId] = lnks
	}
	links := make([]Link, 0, len(lnks))
	for _, val := range lnks {
		links = append(links, val)
	}
	return links, nil
}
