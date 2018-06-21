package capricornius


import (
	"log"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"time"
)

type Partner struct {
	ID    string
	Name  string
	Email string
	Phone int
	Password
	Token
}

type Partners struct{
	Records	[]*Partner
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
	tableCreationPartner := sq.Insert("Partner").Columns("Name","Email","Phone","Password","Token")

	tableGetPartner := sq.Select("*").From("Partner")

	tableUpdatePartner := sq.Update("?").From(Partner)


	
)

func (p *Partner) GetPartnerfromID(db *sql.DB) error {
	FromID := tableGetPartner.Where("ID= ?").ToSql()

	return FromID.QueryRow(p.ID).Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
}

func (p *Partner) GetPartnerfromEmail(db *sql.DB) error {
	FromEmail:= tableGetPartner.Where("Email= ?").ToSql()	
	rows, err := FromEmail.Query(p.email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	partners:= make(Partners, 0)
	for rows.Next(){
		err:= rows.Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
		if err != nil {
			log.Fatal(err)
		}
		partners = append(partners, &Partner)
		
	}
	return partners

}



func (p *Partner) GetPartnerfromName(db *sql.DB)error {
	FromName:= tableGetPartner.Where("Name = ?").ToSql()

	return FromName.QueryRow(p.name).Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
}

func (p. *Partner) GetPartnerfromToken(db *sql.DB) error {
	FromToken:= tableGetPartner.Where("Token = ?").ToSql()

	return FromToken.QueryRow(p.token).Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)

}

func (p *Partner) CreatePartner(db *sql.DB) error {
	values:=tableCreationPartner.Values(p.Name, p.email, p.phone, p.password, p.token)
	sql,args, err := values.ToSql()
	return sql.Exec()
}

func (p *Partner) GetallPartner(db *sql.DB, count int) error {
	all:= tableGetPartner.ToSql()
	rows, err := all.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	partners:= make(Partners, 0)
	for rows.Next(){
		err:= rows.Scan(&p.Name, &p.email, &p.phone, &p.password, &p.token)
		if err != nil {
			log.Fatal(err)
		}
		partners = append(partners, &Partner)
		
	}
	return partners
}

func (p *Partner) UpdatePartner(db *sql.DB) error {
	update:= 
}

func (p *Partner) DeletePartner(db *sql.DB) error {
	tableDeletePartner := sq.Update("Partner").Set("Delete_at": time.date).Where("ID":p.id)
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
