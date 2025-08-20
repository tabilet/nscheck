package vaultcheck

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/openbao/openbao/api/v2"
)

const (
	RootTokenAddr = ".vault-token"
	sleeping      = 4 * time.Second
)

func getClient() (*api.Client, error) {
	time.Sleep(sleeping) // Ensure Namespace is ready before creating a client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	client.SetAddress(os.Getenv("VAULT_ADDR"))
	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))

	token := os.Getenv("VAULT_TOKEN")
	fn := os.Getenv("HOME") + "/" + RootTokenAddr
	if _, err := os.Stat(fn); err == nil {
		bs, err := os.ReadFile(fn)
		if err != nil {
			log.Fatalf("Failed to read root token file: %v", err)
		}
		token = string(bs)
	}
	client.SetToken(token)

	return client, nil
}

func cloneClient(ctx context.Context, client *api.Client, pname string) (*api.Client, error) {
	_, err := client.Logical().WriteWithContext(ctx, "sys/namespaces/"+pname, nil)
	if err != nil {
		return nil, err
	}
	clone, err := client.Clone()
	if err != nil {
		return nil, err
	}
	clone.SetToken(client.Token())
	top := client.Namespace()
	if top == "" {
		clone.SetNamespace(pname)
	} else {
		clone.SetNamespace(top + "/" + pname)
	}
	time.Sleep(sleeping)
	return clone, nil
}

func combinedPath(rootNS string) string {
	ns := os.Getenv("VAULT_NAMESPACE")
	if ns == "" {
		return rootNS
	}
	return ns + "/" + rootNS
}
