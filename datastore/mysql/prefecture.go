package mysql

import (
	"github.com/usk81/minimum-wage-ja/domain"

	"github.com/jmoiron/sqlx"
)

type PrefectureRepository struct {
	DB *sqlx.DB
}

// FindAll gets all of prefectures on database
func (p *PrefectureRepository) FindAll() (result []domain.Prefecture, err error) {

	err = p.DB.Select(&result, query("SELECT id, name FROM `%s`", p.table()))
	return
}

func (p *PrefectureRepository) table() string {
	return "prefectures"
}
