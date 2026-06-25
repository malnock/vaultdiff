package vault

import (
	"os"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	_, err := NewClient("dev", "", "sometoken")
	if err == nil {
		t.Fatal("expected error when VAULT_ADDR is not set")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")

	_, err := NewClient("dev", "http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error when token is not set")
	}
}

func TestNewClient_ExplicitAddressAndToken(t *testing.T) {
	client, err := NewClient("staging", "http://127.0.0.1:8200", "root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Env != "staging" {
		t.Errorf("expected env 'staging', got '%s'", client.Env)
	}
}

func TestNewClient_EnvVarFallback(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "root")

	client, err := NewClient("prod", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.api == nil {
		t.Error("expected non-nil vault api client")
	}
}
