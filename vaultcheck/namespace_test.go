package vaultcheck

import (
	"testing"
)

// TestNamespaceRoot tests the namespace functionality at the root namespace.
func TestNamespaceRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckNamespace(client)
	if err != nil {
		t.Fatalf("Namespace failed: %v", err)
	}
}
