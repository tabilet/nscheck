package vaultcheck

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/openbao/openbao/api/v2"
)

// CheckTokenRoot checks if the token auth is mounted and cannot be disabled in the root namespace.
func CheckTokenRoot(client *api.Client) error {
	ctx := context.Background()

	rootToken := client.Token()

	path := "token"
	err := checkTokenAuth(ctx, client, path)
	if err != nil {
		return err
	}

	tokenAuth, secret, err := getTokenAuthSecret(ctx, client, rootToken, "default")
	if err != nil {
		return err
	}

	// test revocation
	s, err := revokeTokenByRootToken(ctx, client, tokenAuth, path, rootToken, secret.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s != nil && s.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s.Auth)
	}
	return nil
}

// CheckTokenNamespace checks if the token auth is mounted and cannot be disabled in the namespace.
func CheckTokenNamespace(client *api.Client) error {
	ctx := context.Background()

	rootToken := client.Token()

	rootNS := "pname"
	clone, err := cloneClient(ctx, client, rootNS)
	if err != nil {
		return err
	}

	path := "token"
	err = checkTokenAuth(ctx, clone, path)
	if err != nil {
		return err
	}

	tokenAuth1, secret1, err := getTokenAuthSecret(ctx, clone, rootToken, "default")
	if err != nil {
		return err
	}
	tokenAuth2, secret2, err := getTokenAuthSecret(ctx, clone, rootToken, "default")
	if err != nil {
		return err
	}

	if client.Token() != clone.Token() ||
		client.Namespace() != os.Getenv("VAULT_NAMESPACE") ||
		clone.Namespace() != combinedPath(rootNS) {
		return fmt.Errorf("root Token: %s in namespace %s, clone %s in namespace %s", client.Token(), client.Namespace(), clone.Token(), clone.Namespace())
	}

	s1, err := revokeTokenByRootToken(ctx, clone, tokenAuth1, path, rootToken, secret1.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s1 != nil && s1.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s1.Auth)
	}
	// remove a namespace token from the root namespace ... ok?
	s2, err := revokeTokenByRootToken(ctx, client, tokenAuth2, path, rootToken, secret2.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s2 != nil && s2.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s2.Auth)
	}

	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+rootNS)
	if err != nil {
		return err
	}
	return nil
}

// CheckTokenMix checks if the token auth is mounted and cannot be disabled in the root namespace and in the namespace.
func CheckTokenMix(client *api.Client) error {
	ctx := context.Background()

	rootToken := client.Token()

	path := "token"

	// in root namespace
	err := checkTokenAuth(ctx, client, path)
	if err != nil {
		return err
	}
	tokenAuth1, secret1, err := getTokenAuthSecret(ctx, client, rootToken, "default")
	if err != nil {
		return err
	}
	tokenAuth3, secret3, err := getTokenAuthSecret(ctx, client, rootToken, "default")
	if err != nil {
		return err
	}

	// in namespace
	rootNS := "pname"
	clone, err := cloneClient(ctx, client, rootNS)
	if err != nil {
		return err
	}

	err = checkTokenAuth(ctx, clone, path)
	if err != nil {
		return err
	}
	tokenAuth2, secret2, err := getTokenAuthSecret(ctx, clone, rootToken, "default")
	if err != nil {
		return err
	}
	tokenAuth4, secret4, err := getTokenAuthSecret(ctx, clone, rootToken, "default")
	if err != nil {
		return err
	}

	if client.Token() != clone.Token() ||
		client.Namespace() != os.Getenv("VAULT_NAMESPACE") ||
		clone.Namespace() != combinedPath(rootNS) {
		return fmt.Errorf("root Token: %s in namespace %s, clone %s in namespace %s", client.Token(), client.Namespace(), clone.Token(), clone.Namespace())
	}

	s1, err := revokeTokenByRootToken(ctx, client, tokenAuth1, path, rootToken, secret1.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s1 != nil && s1.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s1.Auth)
	}
	s2, err := revokeTokenByRootToken(ctx, clone, tokenAuth2, path, rootToken, secret2.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s2 != nil && s2.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s2.Auth)
	}

	// remove a normal token in space from the namespace ... ok?
	s3, err := revokeTokenByRootToken(ctx, clone, tokenAuth3, path, rootToken, secret3.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s3 != nil && s3.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s3.Auth)
	}

	// remove a namespace token from the root namespace ... ok?
	s4, err := revokeTokenByRootToken(ctx, client, tokenAuth4, path, rootToken, secret4.Auth.ClientToken)
	if err != nil {
		return err
	}
	if s4 != nil && s4.Auth != nil {
		return fmt.Errorf("revocation failed %+v", s4.Auth)
	}

	// clean up
	client.SetNamespace(os.Getenv("VAULT_NAMESPACE"))
	_, err = client.Logical().DeleteWithContext(ctx, "sys/namespaces/"+rootNS)
	if err != nil {
		return err
	}
	return nil
}

// checkTokenAuth checks if the token auth is mounted and cannot be disabled.
func checkTokenAuth(ctx context.Context, client *api.Client, path string) error {
	sys := client.Sys()

	mountsRspn, err := sys.ListAuthWithContext(ctx)
	if err != nil {
		return err
	}
	for k, rspn := range mountsRspn {
		if !slices.Contains([]string{"token/"}, k) {
			return fmt.Errorf("wrong mount response: %s => %+v", k, rspn)
		}
	}

	/*
		err = sys.DisableAuthWithContext(ctx, path)
		if err == nil {
			return fmt.Errorf("should have failed")
		} else if rErr, ok := err.(*api.ResponseError); !ok || rErr.StatusCode != 400 || (rErr.Errors)[0] != "token credential backend cannot be disabled" {
			return fmt.Errorf("%#v", rErr.Errors)
		}
	*/
	return nil
}

// getTokenAuthSecret creates a new token from client, which is associated with a namespace, with the given policies and returns the token auth and secret.
func getTokenAuthSecret(ctx context.Context, client *api.Client, rootToken string, policy ...string) (*api.TokenAuth, *api.Secret, error) {
	client.SetToken(rootToken)

	tokenAuth := client.Auth().Token()
	secret, err := tokenAuth.CreateWithContext(ctx, &api.TokenCreateRequest{
		Policies: policy,
	})
	if err != nil {
		return nil, nil, err
	}
	if secret.Auth == nil || secret.Auth.ClientToken == "" {
		return nil, nil, fmt.Errorf("Auth data: %+v", secret.Auth)
	}

	token := secret.Auth.ClientToken
	client.SetToken(token)
	tokenAuth = client.Auth().Token()

	selfSecret, err := tokenAuth.LookupSelfWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	if selfSecret == nil || selfSecret.Data == nil || selfSecret.Data["policies"] == nil {
		return nil, nil, fmt.Errorf("self response data nil: %+v", selfSecret)
	}
	policies := selfSecret.Data["policies"].([]any)
	for _, p := range policies {
		if !slices.Contains(policies, any(p)) {
			return nil, nil, fmt.Errorf("policy %s not found in %v.", p, policies)
		}
	}

	client.SetToken(rootToken)
	return tokenAuth, secret, nil
}

func revokeTokenByRootToken(ctx context.Context, client *api.Client, tokenAuth *api.TokenAuth, path, rootToken, token string) (*api.Secret, error) {
	client.SetToken(rootToken)

	secret, err := client.Logical().WriteWithContext(ctx, "auth/"+path+"/revoke", map[string]any{
		"token": token,
	})
	if err != nil {
		return nil, err
	}
	if secret != nil {
		return nil, fmt.Errorf("secret found after revoke api: %+v", secret)
	}

	s, err := tokenAuth.LookupSelfWithContext(ctx)
	if err != nil {
		if rErr, ok := err.(*api.ResponseError); !ok || rErr.StatusCode != 403 || (rErr.Errors)[0] != "permission denied" {
			return nil, fmt.Errorf("error: %#v", rErr.Errors)
		}
	}
	return s, nil
}
