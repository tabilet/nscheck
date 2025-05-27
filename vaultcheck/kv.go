package vaultcheck

import (
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/hashicorp/vault/api"
)

// CheckKVRoot checks if the KV secret engine is mounted and can be deleted in the root namespace.
func CheckKVRoot(client *api.Client) error {
	ctx := context.Background()

	path := "secret-v2"
	err := checkKVMount(ctx, client, path)
	if err != nil {
		return err
	}

	kvSecret, err := createGetDeleteKV2(ctx, client, path, "mysecret", "myadmin", "123456")
	if err != nil {
		return err
	}
	if kvSecret != nil && kvSecret.Data != nil {
		return fmt.Errorf("KV secret: %#v", kvSecret.Data)
	}

	kvSecret, err = createGetDeleteKV2(ctx, client, path, "mysecret", "myadmin7", "123456")
	if err != nil {
		return err
	}
	if kvSecret != nil && kvSecret.Data != nil {
		return fmt.Errorf("KV secret: %#v", kvSecret.Data)
	}

	sys := client.Sys()
	err = sys.UnmountWithContext(ctx, path)
	if err != nil {
		return err
	}

	err = checkKVMount(ctx, client, path)
	if err != nil {
		return err
	}

	err = sys.UnmountWithContext(ctx, path)
	if err != nil {
		return err
	}
	return nil
}

// CheckKVNamespace checks if the KV secret engine is mounted and can be deleted in the namespace.
func CheckKVNamespace(client *api.Client) error {
	ctx := context.Background()

	pname := "pname"
	clone, err := cloneClient(ctx, client, pname)
	if err != nil {
		return err
	}

	path := "secret-v2"
	err = checkKVMount(ctx, clone, path)
	if err != nil {
		return err
	}

	kvSecret, err := createGetDeleteKV2(ctx, clone, path, "yoursecret", "myadmin", "123456")
	if err != nil {
		return err
	}
	if kvSecret != nil && kvSecret.Data != nil {
		return fmt.Errorf("KV secret: %#v", kvSecret.Data)
	}

	sys := clone.Sys()
	err = sys.UnmountWithContext(ctx, path)
	if err != nil {
		return err
	}
	err = checkKVMount(ctx, clone, path)
	if err != nil {
		return err
	}
	err = sys.UnmountWithContext(ctx, path)
	if err != nil {
		return err
	}

	err = checkKVMount(ctx, clone, path)
	if err != nil {
		return err
	}

	err = sys.UnmountWithContext(ctx, path)
	if err != nil {
		return err
	}

	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+pname)
	if err != nil {
		return err
	}
	return nil
}

// CheckKVMix checks if the KV secret engine is mounted and can be deleted in the root namespace and in the namespace.
func CheckKVMix(client *api.Client) error {
	ctx := context.Background()

	path := "secret-v2"
	err := checkKVMount(ctx, client, path)
	if err != nil {
		return err
	}

	name1 := "mysecret"
	kv1, err := createGetKV2(ctx, client, path, name1, "myadmin", "123456")
	if err != nil {
		return err
	}

	pname := "pname"
	clone, err := cloneClient(ctx, client, pname)
	if err != nil {
		return err
	}

	err = checkKVMount(ctx, clone, path)
	if err != nil {
		return err
	}

	name2 := "yoursecret"
	kv2, err := createGetKV2(ctx, clone, path, name2, "myadmin", "000000")
	if err != nil {
		return err
	}

	kvSecret, err := kv2.Get(ctx, name1)
	// kvNS tries to get a secret in root namespace
	if err == nil || (err.Error())[:16] != "secret not found" {
		return fmt.Errorf("KV secret: %+v", kvSecret)
	}

	// in root namespace
	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
	kvSecret, err = kv1.Get(ctx, name2)
	// kv2 tries to get a secret in child namespace
	if err == nil || (err.Error())[:16] != "secret not found" {
		return fmt.Errorf("KV secret: %+v", kvSecret)
	}

	// cleanup
	if err = kv2.Delete(ctx, name2); err == nil {
		err = clone.Sys().UnmountWithContext(ctx, path)
	}
	if err != nil {
		return err
	}

	if err = kv1.Delete(ctx, name1); err == nil {
		err = client.Sys().UnmountWithContext(ctx, path)
	}
	if err != nil {
		return err
	}

	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+pname)
	if err != nil {
		return err
	}
	return nil
}

// checkKVMount mounts the KV secret engine at the given path and checks if it is mounted correctly.
func checkKVMount(ctx context.Context, client *api.Client, path string) error {
	sys := client.Sys()

	err := sys.MountWithContext(ctx, path, &api.MountInput{
		Type: "kv-v2",
		Options: map[string]string{
			"upgrade": "false",
		},
	})
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)

	mountsRspn, err := sys.ListMountsWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"secret/", "cubbyhole/", "identity/", "sys/", path + "/"}, k) {
			return fmt.Errorf("mount response: %s => %+v", k, rspn)
		}
	}
	return nil
}

// createGetKV2 creates a KV secret, retrieves it, and confirms that the data is correct.
func createGetKV2(ctx context.Context, client *api.Client, path, name, username, password string) (*api.KVv2, error) {
	kv2 := client.KVv2(path)

	kvSecret, err := kv2.Put(ctx, name, map[string]any{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	if kvSecret.Data != nil {
		return nil, fmt.Errorf("KV secret: %#v", kvSecret.Data)
	}

	kvSecret, err = kv2.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if kvSecret.Data == nil ||
		kvSecret.Data["username"].(string) != username ||
		kvSecret.Data["password"].(string) != password {
		return nil, fmt.Errorf("KV secret: %#v", kvSecret.Data)
	}

	return kv2, nil
}

func createGetDeleteKV2(ctx context.Context, client *api.Client, path, name, username, password string) (*api.KVSecret, error) {
	kv2, err := createGetKV2(ctx, client, path, name, username, password)
	if err != nil {
		return nil, err
	}

	err = kv2.Delete(ctx, name)
	if err != nil {
		return nil, err
	}
	kvSecret, err := kv2.Get(ctx, name)
	if err != nil {
		if rErr, ok := err.(*api.ResponseError); !ok || rErr.StatusCode != 404 || (rErr.Errors)[0] != "not found" {
			return nil, fmt.Errorf("error: %#v", rErr.Errors)
		}
	}
	return kvSecret, nil
}
