package mysql

import (
	"database/sql"
	"time"

	"github.com/usk81/minimum-wage-ja/domain"
	"github.com/usk81/minimum-wage-ja/interface/repository"

	"github.com/jmoiron/sqlx"
)

const querySetWages = `INSERT INTO %s 
(
	id,
	prefecture_id,
	hourly,
	daily,
	name,
	regional,
	implemented_at
)
VALUE (
	:id,
	:prefecture_id,
	:hourly,
	:daily,
	:name,
	:regional,
	:implemented_at
)
ON DUPLICATE KEY UPDATE 
hourly = :hourly,
daily = :daily,
implemented_at = :implemented_at;`

const querySetWageLogs = `INSERT INTO %s
(
	id,
	prefecture_id,
	hourly,
	daily,
	name,
	regional,
	implemented_at
)
VALUE (
	:id,
	:prefecture_id,
	:hourly,
	:daily,
	:name,
	:regional,
	:implemented_at
)`

const queryWageFind = `SELECT
	id,
	prefecture_id,
	hourly,
	daily,
	name,
	regional,
	implemented_at
FROM
	%s
WHERE name = ?
AND prefecture_id = ?
AND implemented_at <= ?
ORDER BY updated_at DESC
LIMIT 1`

const queryWageFindByID = `SELECT
	id,
	prefecture_id,
	hourly,
	daily,
	name,
	regional,
	implemented_at
FROM
	%s
WHERE id = ?
AND implemented_at <= ?
ORDER BY updated_at DESC
LIMIT 1`

// FIXME:
const queryWageFindByPrefectureID = `SELECT
	id,
	prefecture_id,
	hourly,
	daily,
	name,
	regional,
	implemented_at
FROM
	%s
WHERE prefecture_id = ?
AND implemented_at <= ?
ORDER BY id, updated_at DESC
GROUP BY id`

type WageRepository struct {
	DB *sqlx.DB
}

func (w *WageRepository) Find(name string, prefectureID int) (result *domain.Wage, err error) {
	result = &domain.Wage{}
	err = w.DB.Get(result, w.DB.Rebind(query(queryWageFind, w.table())), name, prefectureID, time.Now())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (w *WageRepository) FindByID(id int) (result *domain.Wage, err error) {
	result = &domain.Wage{}
	err = w.DB.Get(result, w.DB.Rebind(query(queryWageFindByID, w.table())), id, time.Now())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (w *WageRepository) FindByPrefectureID(id int) (result []domain.Wage, err error) {
	err = w.DB.Select(result, w.DB.Rebind(query(queryWageFindByPrefectureID, w.table())), id, time.Now())
	return
}

func (w *WageRepository) Set(req domain.Wage) (err error) {
	tx, err := w.DB.Beginx()
	if err != nil {
		return
	}
	if _, err = tx.NamedExec(query(querySetWageLogs, w.table()), req); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (w *WageRepository) table() string {
	return "wage_logs"
}

func (w *WageRepository) Checksum() repository.Checksum {
	return &ChecksumRepository{
		DB: w.DB,
	}
}

func (w *WageRepository) Prefecture() repository.Prefecture {
	return &PrefectureRepository{
		DB: w.DB,
	}
}
