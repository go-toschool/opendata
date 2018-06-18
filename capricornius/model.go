package main

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"

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
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableCreationPartner = sq.Insert("Partner").Columns("Name","Email","Phone","Password","Token")

	tableGetPartner = sq.Select("*").From("Partner")

	
)

func (p *Partner) GetPartnerforID(db *sql.DB) error {
	ForID := tableGetPartner.Where("ID":p.ID).ToSql()

	return db.QueryRow(ForID).Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
}

func (p *Partner) GetPartnerforEmail(db *sql.DB) error {
	ForEmail:= tableGetPartner.Where("Email": p.Email).ToSql()	

	return db.QueryRow(ForEmail).Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
}

func (p *Partner) GetPartnerforName(db *sql.DB)error {
	ForName:= tableGetPartner.Where("Name" : p.Name).ToSql()

	return db.QueryRow(ForName).(&p.Name, &p.email, &p.phone, &p.password, &p.token)
}

func (p. *Partner) GetPartnerforToken(db *sql.DB) error {
	ForToken:= tableGetPartner.Where
}

func (p *Partner) CreatePartner(db *sql.DB) error {
	values:=tableCreationPartner.Values(p.Name, p.email, p.phone, p.password, p.token)
	sql,args, err := values.ToSql()
	return db.QueryRow(sql)
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
