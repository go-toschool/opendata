package capricornius

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Partner struct {
	ID       string
	Name     string
	Email    string
	Phone    int
	Password string
	Token    string
}

type Callback struct {
	UserID     string
	ID         string
	URL        string
	HTTPMethod string
	CreatedAt  time.Time
	DeleteAt   time.Time
}

type Statistis struct {
	UserID      string
	CallbaksID  string
	Request     string
	augResponse string
}

func Init() {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableCreationCallback := psql.Insert("Callback").Columns("User_id", "ID", "URL", "HTTP_method", "Created_at")

	tableGetPartner := psql.Select("*").From("Partner")
	tableCreationPartner := psql.Insert("Partner").Columns("Name", "Email", "Phone", "Password", "Token")

	tableGetCallback := psql.Select("*").From("Calback")
}

func (p *Partner) getPartnerfromID(db *sql.DB) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tableGetPartner := psql.Select("*").From("Partner")
	FromID := tableGetPartner.Where(sq.Eq{"ID": p.ID})
	getPartner := FromID.ToSql()

	return getPartner.QueryRow(p.ID).Scan(&p.Name, &p.Email, &p.Phone, &p.Password, &p.Token)
}

func (p *Partner) getPartnerfromEmail(db *sql.DB) ([]*Partner, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableGetPartner := psql.Select("*").From("Partner")
	FromEmail := tableGetPartner.Where("Email = ?").ToSql()
	rows, err := FromEmail.Query(p.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	partners := make([]*Partner, 0)
	for rows.Next() {
		err := rows.Scan(&p.Name, &p.Email, &p.Phone, &p.Password, &p.Token)
		if err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}

	return partners, nil
}

func (p *Partner) getPartnerfromName(db *sql.DB) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableGetPartner := psql.Select("*").From("Partner")
	FromName := tableGetPartner.Where("Name = ?", p.Name).ToSql()

	return FromName.QueryRow(FromName).Scan(&p.Name, &p.Email, &p.Phone, &p.Password, &p.Token)
}

func (p *Partner) getPartnerfromToken(db *sql.DB) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableGetPartner := psql.Select("*").From("Partner")
	FromToken := tableGetPartner.Where("Token = ?", p.Token).ToSql()
	return FromToken.QueryRow(FromToken).Scan(&p.Name, &p.Email, &p.Phone, &p.Password, &p.Token)

}

func (p *Partner) createPartner(db *sql.DB) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableCreationPartner := psql.Insert("Partner").Columns("Name", "Email", "Phone", "Password", "Token")

	values := tableCreationPartner.Values(p.Name, p.Email, p.Phone, p.Password, p.Token)
	sql, args, err := values.ToSql()
	_, error := db.Exec(sql)
	return error
}

func (p *Partner) getallPartner(db *sql.DB, count int) ([]*Partner, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableGetPartner := psql.Select("*").From("Partner")
	all := tableGetPartner.ToSql()
	rows, err := db.Query(all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	partners := make([]*Partner, 0)
	for rows.Next() {
		err := rows.Scan(&p.Name, &p.Email, &p.Phone, &p.Password, &p.Token)
		if err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}
	return partners, nil
}

func (p *Partner) UpdatePartner(db *sql.DB) error {
	tableUpdatePartner := sq.Update("Partner").Where(sq.Eq{"Id": &p.ID})
	switch {
	case &p.Name != nil:
		tableUpdatePartner = tableUpdatePartner.Set("Name", &p.Name)

	case &p.Email != nil:
		tableUpdatePartner = tableUpdatePartner.Set("Email", &p.Email)

	case &p.Phone != nil:
		tableUpdatePartner = tableUpdatePartner.Set("Phone", &p.Phone)

	case &p.Password != nil:
		tableUpdatePartner = tableUpdatePartner.Set("Password", &p.Password)

	case &p.Token != nil:
		tableUpdatePartner = tableUpdatePartner.Set("Token", &p.Token)
	}

	sql, args, err := tableUpdatePartner.ToSql()
	db.Exec(sql)
	return err
}

func (c *Callback) createCallback(db *sql.DB) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableCreationCallback := psql.Insert("Callback").Columns("User_id", "ID", "URL", "HTTP_method", "Created_at")

	values := tableCreationCallback.Values(c.UserID, c.ID, c.URL, c.HTTPMethod, time.Now)
	sql, args, err := values.ToSql()
	_, error := db.Exec(sql)
	return error
}

func (c *Callback) updateCallback(db *sql.DB) error {
	tableUpdateCallback := sq.Update("Callback").Where(sq.Eq{"Id": &c.ID})
	switch {
	case &c.URL != nil:
		tableUpdateCallback = tableUpdateCallback.Set("URL", &c.URL)

	case &c.HTTPMethod != nil:
		tableUpdateCallback = tableUpdateCallback.Set("HTTP_method", &c.HTTPMethod)
	}

	sql, args, err := tableUpdateCallback.ToSql()
	db.Exec(sql)
	return err
}

func (c *Callback) DeleteCallback(db *sql.DB) error {
	tableDeleteCallback := psql.Update("Callback").Where("ID= ?", &c.ID).Set("Delete_at= ?", time.Now).ToSql()
	_, err := db.Exec(tableDeleteCallback)
	return err
}

func (c *Callback) GetCallbackfromID(db *sql.DB) error {
	fromID := tableGetCallback.Where("ID = ?", c.ID).Tosql()
	return db.QueryRow(fromID).Scan(&c.ID, &c.UserID, &c.URL, &c.HTTPMethod, &c.CreatedAt, &c.DeleteAt)
}
