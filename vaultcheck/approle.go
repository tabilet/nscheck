package vaultcheck

import (
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

// CheckApproleRoot checks if the AppRole auth is mounted and can be deleted in the root namespace.
func CheckApproleRoot(client *api.Client) error {
	ctx := context.Background()

	sys := client.Sys()

	path := "approle"
	err := sys.EnableAuthWithOptionsWithContext(ctx, path, &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		return err
	}
	mountsRspn, err := sys.ListAuthWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"token/", "approle/"}, k) {
			return fmt.Errorf("mount response: %s => %+v", k, rspn)
		}
	}

	err = sys.DisableAuthWithContext(ctx, path)
	if err != nil {
		return err
	}
	mountsRspn, err = sys.ListAuthWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"token/"}, k) {
			return fmt.Errorf("mount response: %s => %+v", k, rspn)
		}
	}

	_, secretID, clientToken, err := getApprole(client, ctx, path, "myrole")
	if err != nil {
		return err
	}
	if clientToken == "" {
		return fmt.Errorf("no client token")
	}

	err = dropApprole(client, ctx, secretID, path, "myrole")
	if err != nil {
		return err
	}

	return nil
}

// CheckApproleNamespace checks if the AppRole auth is mounted and can be deleted in the namespace.
func CheckApproleNamespace(client *api.Client) error {
	ctx := context.Background()

	rootNS := "pname"
	clone, err := cloneClient(ctx, client, rootNS)
	if err != nil {
		return err
	}

	sys := clone.Sys()

	path := "approle"
	err = sys.EnableAuthWithOptionsWithContext(ctx, path, &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		return err
	}
	mountsRspn, err := sys.ListAuthWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"token/", "approle/"}, k) {
			return fmt.Errorf("mount response: %s => %+v", k, rspn)
		}
	}

	err = sys.DisableAuthWithContext(ctx, path)
	if err != nil {
		return err
	}
	mountsRspn, err = sys.ListAuthWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"token/"}, k) {
			return fmt.Errorf("mount response: %s => %+v", k, rspn)
		}
	}

	_, secretID, clientToken, err := getApprole(clone, ctx, path, "myrole")
	if err != nil {
		return err
	}
	if clientToken == "" {
		return fmt.Errorf("no client token")
	}

	err = dropApprole(clone, ctx, secretID, path, "myrole")
	if err != nil {
		return err
	}

	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+rootNS)
	if err != nil {
		return err
	}

	return nil
}

// CheckApproleNamespace checks if the AppRole auth is mounted and can be deleted in the namespace.
func CheckApproleMix(client *api.Client) error {
	ctx := context.Background()

	rootNS := "pname"
	clone, err := cloneClient(ctx, client, rootNS)
	if err != nil {
		return err
	}

	path := "approle"

	var secret *api.Secret
	var roleID, secretID, clientToken, roleNS, secretNS, clientTokenNS string

	roleID, secretID, clientToken, err = getApprole(client, ctx, path, "myrole")
	if err != nil {
		return err
	}
	if clientToken == "" {
		return fmt.Errorf("no client token")
	}

	roleNS, secretNS, clientTokenNS, err = getApprole(clone, ctx, path, "yourrole")
	if err != nil {
		return err
	}
	if clientTokenNS == "" {
		return fmt.Errorf("no client token")
	}

	auth, err := approle.NewAppRoleAuth(roleID, &approle.SecretID{FromString: secretID})
	if err != nil {
		return err
	}
	secret, err = auth.Login(ctx, clone)
	if err == nil {
		return fmt.Errorf("error should exist, but we got nil. secret id: %#v", secret)
	}

	authNS, err := approle.NewAppRoleAuth(roleNS, &approle.SecretID{FromString: secretNS})
	if err != nil {
		return err
	}
	secret, err = authNS.Login(ctx, client)
	if err == nil {
		return fmt.Errorf("error should exist, but we got nil. secret id: %#v", secret)
	}

	err = dropApprole(clone, ctx, secretNS, path, "yourrole")
	if err != nil {
		return err
	}

	err = dropApprole(client, ctx, secretID, path, "myrole")
	if err != nil {
		return err
	}

	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+rootNS)
	if err != nil {
		return err
	}
	return nil
}

func getApprole(client *api.Client, ctx context.Context, path, roleName string, policies ...string) (roleID, secretID, token string, err error) {
	if len(policies) == 0 {
		policies = []string{"default"}
	}

	sys := client.Sys()
	logical := client.Logical()

	err = sys.EnableAuthWithOptionsWithContext(ctx, path, &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		return "", "", "", err
	}
	time.Sleep(time.Second * 2)

	_, err = logical.WriteWithContext(ctx, "auth/"+path+"/role/"+roleName, map[string]any{
		"policies": policies,
	})
	if err != nil {
		return "", "", "", err
	}
	secret, err := logical.WriteWithContext(ctx, "auth/"+path+"/role/"+roleName+"/secret-id", nil)
	if err != nil {
		return "", "", "", err
	}
	secretID = secret.Data["secret_id"].(string)
	secret, err = logical.ReadWithContext(ctx, "auth/"+path+"/role/"+roleName+"/role-id")
	if err != nil {
		return "", "", "", err
	}
	roleID = secret.Data["role_id"].(string)

	auth, err := approle.NewAppRoleAuth(roleID, &approle.SecretID{FromString: secretID})
	if err != nil {
		return "", "", "", err
	}
	secret, err = auth.Login(ctx, client)
	if err != nil {
		return "", "", "", err
	}
	if secret.Auth == nil {
		err = fmt.Errorf("No auth data")
		return "", "", "", err
	}

	token = secret.Auth.ClientToken

	return roleID, secretID, token, nil
}

func dropApprole(client *api.Client, ctx context.Context, secretID, path, roleName string) error {
	sys := client.Sys()
	logical := client.Logical()

	_, err := logical.WriteWithContext(ctx, "auth/"+path+"/role/"+roleName+"/secret-id/destroy", map[string]any{
		"secret_id": secretID,
	})
	if err != nil {
		return err
	}

	_, err = logical.DeleteWithContext(ctx, "auth/"+path+"/role/"+roleName)
	if err != nil {
		return err
	}

	secret, err := logical.ListWithContext(ctx, "auth/"+path+"/role")
	if err != nil {
		return err
	}
	if secret != nil {
		return fmt.Errorf("List response: %+v", secret)
	}

	return sys.DisableAuthWithContext(ctx, path)
}
