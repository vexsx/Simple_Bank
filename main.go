package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/vexsx/Simple-Bank/api"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable"
	serverAdderss = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("error with opening db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAdderss)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}
