package vaultcheck

import (
	"testing"
)

// TestTokenRoot tests the token authentication method at the root namespace.
func TestTokenRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckTokenRoot(client)
	if err != nil {
		t.Fatalf("KVRoot failed: %v", err)
	}
}

// TestTokenNamespace tests the token authentication method at a namespace.
func TestTokenNamespace(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckTokenNamespace(client)
	if err != nil {
		t.Fatalf("TokenNamespace failed: %v", err)
	}
}

// TestTokenMix tests the token authentication method at both the root namespace and a namespace.
func TestTokenMix(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckTokenMix(client)
	if err != nil {
		t.Fatalf("TokenMix failed: %v", err)
	}
}
