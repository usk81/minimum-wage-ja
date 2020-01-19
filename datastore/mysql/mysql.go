package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/usk81/minimum-wage-ja/interface/repository"
)

func query(q, table string) string {
	a := fmt.Sprintf(q, table)
	// fmt.Println(a)
	return a
}

type Connector struct {
	SCHEME   string `envconfig:"scheme" default:"mysql"`
	User     string `envconfig:"user" required:"true"`
	Password string `envconfig:"password" required:"true"`
	IP       string `envconfig:"ip" required:"true"`
	Port     string `envconfig:"port" default:"3306"`
	DATABASE string `envconfig:"database" required:"true"`
}

func (c *Connector) Connect(ctx context.Context) (repo repository.Wage, err error) {
	var db *sql.DB
	if c.SCHEME == "mysql" {
		urlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", c.User, c.Password, c.IP, c.Port, c.DATABASE)
		if db, err = sql.Open(c.SCHEME, urlstr); err != nil {
			return
		}
	}
	return &WageRepository{
		DB: sqlx.NewDb(db, "mysql"),
	}, nil
}
