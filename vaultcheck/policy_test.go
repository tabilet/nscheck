package vaultcheck

import (
	"testing"
)

// TestPolicyRoot tests the policy functionality at the root namespace.
func TestPolicyRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyRootDefault(client)
	if err != nil {
		t.Fatalf("PolicyRootDefault failed: %v", err)
	}
}

// TestPolicyRootCustom tests the policy functionality at the root namespace with custom policies.
func TestPolicyRootCustom(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyRootCustom(client)
	if err != nil {
		t.Fatalf("PolicyRootCustom failed: %v", err)
	}
}

// TestPolicyNamespaceDefault tests the policy functionality in a specific namespace with default policies.
func TestPolicyNamespaceDefault(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyNamespaceDefault(client)
	if err != nil {
		t.Fatalf("PolicyNamespaceDefault failed: %v", err)
	}
}

// TestPolicyNamespaceCustom tests the policy functionality in a specific namespace with custom policies.
func TestPolicyNamespaceCustom(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyNamespaceCustom(client)
	if err != nil {
		t.Fatalf("PolicyNamespaceCustom failed: %v", err)
	}
}

// TestPolicyMixDeleteInNamespace tests the policy deletion functionality in a specific namespace.
func TestPolicyMixDeleteInNamespace(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyMixDeleteInNamespace(client)
	if err != nil {
		t.Fatalf("PolicyMixDeleteInNamespace failed: %v", err)
	}
}

// TestPolicyMixDeleteInRoot tests the policy deletion functionality in the root namespace.
func TestPolicyMixDeleteInRoot(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = CheckPolicyMixDeleteInRoot(client)
	if err != nil {
		t.Fatalf("PolicyMixDeleteInRoot failed: %v", err)
	}
}
