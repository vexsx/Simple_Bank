package db

import (
	"database/sql"
	"github.com/vexsx/Simple-Bank/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("error")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("erro = ", err)

	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
