package domain

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// Prefecture defines a prefecture in Japan
type Prefecture struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Wage struct {
	ID            string      `json:"id" db:"id"`
	PrefectureID  int         `json:"-" db:"prefecture_id"`
	Prefecture    *Prefecture `json:"prefecture" db:"-"`
	Hourly        int         `json:"hourly" db:"hourly"`
	Daily         int         `json:"daily" db:"daily"`
	Name          string      `json:"name" db:"name"`
	Regional      bool        `json:"regional" db:"regional"`
	ImplementedAt time.Time   `json:"implemented_at" db:"implemented_at"`
	DeletedAt     *time.Time  `json:"deleted_at" db:"deleted_at"`
}

// ID creates a new ID
func ID() string {
	tt := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(tt.UnixNano())), 0)
	return strings.ToLower(ulid.MustNew(ulid.Timestamp(tt), entropy).String())
}
