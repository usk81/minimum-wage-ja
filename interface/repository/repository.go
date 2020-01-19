package repository

import (
	"context"

	"github.com/usk81/minimum-wage-ja/domain"
)

type Connector interface {
	Connect(ctx context.Context) (Wage, error)
}

type Wage interface {
	Find(name string, prefectureID int) (result *domain.Wage, err error)
	FindByID(id int) (result *domain.Wage, err error)
	FindByPrefectureID(id int) (result []domain.Wage, err error)
	Set(domain.Wage) (err error)
	Checksum() Checksum
	Prefecture() Prefecture
}

type Checksum interface {
	Get(key string) (cs string, err error)
	Set(key string, cs string) (err error)
}

type Prefecture interface {
	FindAll() (result []domain.Prefecture, err error)
}
