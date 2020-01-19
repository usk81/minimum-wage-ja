package mysql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type ChecksumRepository struct {
	DB *sqlx.DB
}

type tmpChecksum struct {
	Name      string    `db:"name"`
	Checksum  string    `db:"checksum"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *ChecksumRepository) Get(name string) (cs string, err error) {
	var tc tmpChecksum
	if err = c.DB.Get(&tc, query("SELECT `checksum` FROM %s WHERE `name` = ? LIMIT 1", c.table()), name); err != nil && err != sql.ErrNoRows {
		return
	}
	return tc.Checksum, nil
}

func (c *ChecksumRepository) Set(name string, cs string) (err error) {
	q := query("INSERT INTO %s (`name`, `checksum`, `updated_at`) VALUES (:name, :checksum, :updated_at) ON DUPLICATE KEY UPDATE `checksum` = :checksum, updated_at = :updated_at", c.table())
	_, err = c.DB.NamedExec(q, tmpChecksum{
		Name:      name,
		Checksum:  cs,
		UpdatedAt: time.Now(),
	})
	return
}

func (c *ChecksumRepository) table() string {
	return "checksums"
}
