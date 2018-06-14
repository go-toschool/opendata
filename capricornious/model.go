package main

import (
	"database/sql"
)

type Partner struct {
	ID    string
	Name  string
	Email string
	Phone int
	Password
	Token
}

type Callbaks struct {
	User_id     string
	ID          string
	URL         string
	HTTP_method string
	Created_at
	Delete_at
}

type Statistis struct {
	User_id      string
	Callbaks_id  string
	Request      string
	aug_response string
}

const (
	tableCreationPartner = `CREATE TABLE partners (
	id SERIAL  PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL,
		phone INT,
		password ,
		token TEXT

	)`

	tableGetPartner = ` SELECT partners (
		id,
		name,
		email
	)
	`
)

func (p *Partner) GetPartnerforID(db *sql.DB) error {

	return db.QueryRow("SELECT name, email, phone FROM partners WHERE id=$1").Scan(&p.Name, &p.email, &p.phone)
}

func (p *Partner) GetPartnerforEmail(db *sql.DB) error {
	return db.QueryRow("SELECT id,name.phone FROM partners WHERE ")
}

func (p *Partner) CreatePartner(db *sql.DB) error {

	return error.New("not implemented yet")
}

func (p *Partner) GetallPartner(db *sql.DB, count int) error {

	return error.New("not implemented yet")
}

func (p *Partner) UpdatePartner(db *sql.DB) error {

}

func (c *Callback) CreateCallback(db *sql.DB) error {

	return error.New("not implemented yet")
}

func (c *Callback) UpdateCallback(db *sql.DB) error {

}

func (p *Partner) Partnercomprobation(db *sql.DB) error {

}

func (c *Callback) Partnercomprobation(db *sql.DB) error {

}
