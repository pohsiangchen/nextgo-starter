package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func DataSourceName(port int, host, user, password, dbname, sslmode string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func NewSqlx(driverName, dataSourceName string) *sqlx.DB {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
