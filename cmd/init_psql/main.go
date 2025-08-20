// Package main is the main package for the kinet application.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/openbao/openbao/api/v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: init_psql -addr=http://127.0.0.1:8200\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var addr string

func init() {
	flag.Usage = usage
	flag.StringVar(&addr, "addr", "http://127.0.0.1:8200", "Address of the Vault server")
	flag.Parse()
}

func main() {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	sys := client.Sys()
	rspn, err := sys.Init(&api.InitRequest{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	rootToken := rspn.RootToken
	client.SetToken(rootToken)
	fn, err := os.Create(os.Getenv("HOME") + "/.vault-token")
	if err != nil {
		log.Fatal(err)
	}
	defer fn.Close()
	log.Printf("Root token %s is saved to %s", rootToken, fn.Name())

	_, err = fn.WriteString(rootToken)
	if err != nil {
		log.Fatal(err)
	}

	status, err := sys.Unseal(rspn.Keys[0])
	if err != nil {
		log.Fatal(err)
	}

	if status.Sealed {
		log.Fatalf("failed to unseal Vault: %v", status)
	}
}
