package vaultcheck

import (
	"testing"
)

// TestACLRoot is a test function that checks the ACL policies for root tokens.
func TestACLRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckACLRoot(client)
	if err != nil {
		t.Fatalf("ACLRoot failed: %v", err)
	}
}

// TestACLNamespace is a test function that checks the ACL policies for namespace tokens.
func TestACLNamespace(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckACLNamespace(client)
	if err != nil {
		t.Fatalf("ACLNamespace failed: %v", err)
	}
}

// TestACLMixNormal is a test function that checks the ACL policies for normal tokens across default space and namespace.
func TestACLMixNormal(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckACLMixNormal(client)
	if err != nil {
		t.Fatalf("ACLMixNormal failed: %v", err)
	}
}

// TestACLMixPower is a test function that checks the ACL policies for power tokens across default space and namespace, and 1-level downstream namespace.
func TestACLMixPower(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckACLMixPower(client)
	if err != nil {
		t.Fatalf("ACLMixPower failed: %v", err)
	}
}
