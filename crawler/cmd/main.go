package main

import (
	"github.com/usk81/minimum-wage-ja/crawler"
	"github.com/usk81/minimum-wage-ja/registry"
)

func main() {
	conn, err := registry.NewMySQLDB()
	if err != nil {
		panic(err)
	}
	if err = crawler.Run(conn); err != nil {
		panic(err)
	}
}
