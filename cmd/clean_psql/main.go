// Package main is the main package for the kinet application.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: init_psql -connURL=postgres://user123:secret123!@vm0:5432/openbao\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var connURL string

func init() {
	flag.Usage = usage
	flag.StringVar(&connURL, "connURL", "postgres://user123:secret123!@vm0:5432/openbao", "Address of the PostgreSQL server")
	flag.Parse()
}

func main() {
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	unquoted_table := "openbao_kv_store"
	_, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", unquoted_table))
	if err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
}
