package vaultcheck

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/openbao/openbao/api/v2"
)

// CheckNamespace checks if the namespaces are created and can be deleted.
func CheckNamespace(client *api.Client) error {
	ctx := context.Background()

	logical := client.Logical()

	rootNS := os.Getenv("VAULT_NAMESPACE")
	for _, ns := range []string{"pname", "cname", "dname", "ename"} {
		client.SetNamespace(rootNS)
		_, err := logical.WriteWithContext(ctx, "sys/namespaces/"+ns, nil)
		if err != nil {
			return err
		}
		rspn, err := logical.ListWithContext(ctx, "sys/namespaces")
		if err != nil {
			return err
		}
		if rspn.Data == nil ||
			rspn.Data["keys"] == nil ||
			!slices.Contains(rspn.Data["keys"].([]any), any(ns+"/")) {
			return fmt.Errorf("Namespace list of %s: %+v", ns, rspn.Data)
		}
		rootNS += "/" + ns
	}

	client.SetNamespace(combinedPath("pname/cname"))
	_, err := logical.DeleteWithContext(ctx, "sys/namespaces/dname")
	if err == nil {
		return fmt.Errorf("Delete dname when ename exists: %s", err)
	}

	for _, ns := range []string{"ename", "dname", "cname", "pname"} {
		rootNS = strings.TrimSuffix(rootNS, "/"+ns)
		client.SetNamespace(rootNS)
		_, err := logical.DeleteWithContext(ctx, "sys/namespaces/"+ns)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 4)
		rspn, err := logical.ListWithContext(ctx, "sys/namespaces")
		if err != nil {
			return err
		}
		if rootNS != "" && rspn != nil { // nil is correct response for zero sub-namespace
			return fmt.Errorf("after delete %s, Namespace list of %s => %+v", ns, rootNS, rspn)
		}
	}
	return nil
}
