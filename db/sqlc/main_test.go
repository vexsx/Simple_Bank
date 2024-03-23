package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("erro = ", err)

	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
