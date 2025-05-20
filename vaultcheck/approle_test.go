package vaultcheck

import (
	"testing"
)

// TestApproleRoot tests the AppRole authentication method at the root namespace.
func TestApproleRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckApproleRoot(client)
	if err != nil {
		t.Fatalf("ApproleRoot failed: %v", err)
	}
}

// TestApproleNamespace tests the AppRole authentication method in a specific namespace.
func TestApproleNamespace(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckApproleNamespace(client)
	if err != nil {
		t.Fatalf("ApproleNamespace failed: %v", err)
	}
}

// TestApproleMix tests the AppRole authentication method with mixed namespaces.
func TestApproleMix(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckApproleMix(client)
	if err != nil {
		t.Fatalf("ApproleMix failed: %v", err)
	}
}
