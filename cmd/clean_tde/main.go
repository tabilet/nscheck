// Package main is the main package for the kinet application.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/taosdata/driver-go/v3/taosSql"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: clean_tde -connURL=root:taosdata@tcp(vm0:6030)/openbao\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var connURL string

func init() {
	flag.Usage = usage
	flag.StringVar(&connURL, "connURL", "root:taosdata@tcp(vm0:6030)/openbao", "Address of the TDE server")
	flag.Parse()
}

func main() {
	db, err := sql.Open("taosSql", connURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	for _, table := range []string{"superbao", "supermount"} {
		_, err = db.Exec(fmt.Sprintf("DROP STABLE %s", table))
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}
