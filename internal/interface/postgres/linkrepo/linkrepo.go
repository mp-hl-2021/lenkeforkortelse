package linkrepo

import (
	"database/sql"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/link"
	"github.com/mp-hl-2021/lenkeforkortelse/internal/domain/status"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryGetLinkFieldById = `
	select link from links where linkid = $1
`

func (p *Postgres) CheckIfLinkExists(linkId string) bool {
	l := link.Link{}
	row := p.conn.QueryRow(queryGetLinkFieldById, linkId)
	err := row.Scan(&l.Link, &l.AccountId)
	if err != nil {
		// todo: this function should return (bool, error)
		return false
	}
	return true
}

const queryCreateLink = `
	insert into  links(linkId, link, accountId) VALUES ($1, $2, $3)
	returning linkid
`

func (p *Postgres) StoreLink(lnk link.Link) (link.Link, error) {
	// todo: StoreLink should return just (error)
	accountId := ""
	if lnk.AccountId != nil {
		accountId = *lnk.AccountId
	}
	row := p.conn.QueryRow(queryCreateLink, lnk.LinkId, lnk.Link, accountId)
	tmp := ""
	err := row.Scan(&tmp)
	if err != nil && err == sql.ErrNoRows {
		return lnk, link.ErrAlreadyExist
	}
	return lnk, err
}

const queryDeleteLink = `
	delete from links where linkid = $1
`

func (p *Postgres) DeleteLink(linkId string) error {
	_, err := p.conn.Exec(queryDeleteLink, linkId)
	return err
}

func (p *Postgres) GetLinkByLinkId(linkId string) (link.Link, error) {
	// this function doesn't fill accountId field
	l := link.Link{}
	row := p.conn.QueryRow(queryGetLinkFieldById, linkId)
	err := row.Scan(&l.Link)
	if err != nil && err == sql.ErrNoRows {
		return l, link.ErrNotFound
	}
	return l, err
}

const queryLinksByAccount = `
	select linkId, link, linkStatus from links where accountid = $1
`

func (p *Postgres) GetLinksByAccountId(accountId string) ([]link.Link, error) {
	rows, err := p.conn.Query(queryLinksByAccount, accountId)
	if err != nil {
		return []link.Link{}, err
	}
	defer rows.Close()

	links := make([]link.Link, 0)
	for rows.Next() {
		lnk := link.Link{}
		if err := rows.Scan(&lnk.LinkId, &lnk.Link, &lnk.LinkStatus); err != nil {
			return []link.Link{}, err
		}
		links = append(links, lnk)
	}
	if err := rows.Err(); err != nil {
		return []link.Link{}, err
	}
	return links, nil
}

const queryUpdateLinkStatus = `
	update links
	set linkstatus = $2
	where linkid = $1
`

func (p *Postgres) UpdateLinkStatusByLinkId(linkId string, linkStatus status.LinkStatus) error {
	_, err := p.conn.Exec(queryUpdateLinkStatus, linkId, linkStatus)
	return err
}

const queryGetAllUserLinks = `
	select linkId, link, linkStatus from links where accountid != ''
`

func (p *Postgres) GetAllUserLinks() ([]link.Link, error) {
	rows, err := p.conn.Query(queryGetAllUserLinks)
	if err != nil {
		return []link.Link{}, err
	}
	defer rows.Close()

	links := make([]link.Link, 0)
	for rows.Next() {
		lnk := link.Link{}
		if err := rows.Scan(&lnk.LinkId, &lnk.Link, &lnk.LinkStatus); err != nil {
			return []link.Link{}, err
		}
		links = append(links, lnk)
	}
	if err := rows.Err(); err != nil {
		return []link.Link{}, err
	}
	return links, nil
}
