package registry

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/usk81/minimum-wage-ja/datastore/mysql"
	"github.com/usk81/minimum-wage-ja/interface/repository"
)

func NewMySQLDB() (conn repository.Connector, err error) {
	var mc mysql.Connector
	if err = envconfig.Process("mysql", &mc); err != nil {
		return
	}
	return &mc, nil
}
