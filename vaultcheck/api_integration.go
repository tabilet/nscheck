package vaultcheck

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

func getClient() (*api.Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	client.SetAddress(os.Getenv("VAULT_ADDR"))
	client.SetToken(os.Getenv("VAULT_TOKEN"))
	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
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
	time.Sleep(time.Second * 2)
	return clone, nil
}

func combinedPath(rootNS string) string {
	ns := os.Getenv("VAULT_NAMESPACE")
	if ns == "" {
		return rootNS
	}
	return ns + "/" + rootNS
}
