package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func DataSourceName(port int, host, user, password, dbname, sslmode string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

// Creates a pool of zero or more underlying connections.
func NewDB(driverName, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	Ping(context.Background(), db)

	return db
}

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func Ping(ctx context.Context, db *sql.DB) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
}
