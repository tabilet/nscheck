package vaultcheck

import (
	"testing"
)

// TestKVRoot tests the KV secret engine at the root namespace.
func TestKVRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckKVRoot(client)
	if err != nil {
		t.Fatalf("KVRoot failed: %v", err)
	}
}

// TestKVNamespace tests the KV secret engine at a child namespace.
func TestKVNamespace(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckKVNamespace(client)
	if err != nil {
		t.Fatalf("KVNamespace failed: %v", err)
	}
}

// TestKVMix tests the KV secret engine at both the root and child namespaces.
func TestKVMix(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckKVMix(client)
	if err != nil {
		t.Fatalf("KVMix failed: %v", err)
	}
}
