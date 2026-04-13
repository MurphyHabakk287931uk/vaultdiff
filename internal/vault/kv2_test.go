package vault

import (
	"context"
	"testing"
)

func TestKV2Client_RewritesPath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/data/myapp": {"key": "value"},
	})
	client := NewKV2Client(mock, "secret")

	secrets, err := client.ReadSecrets(context.Background(), "secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Errorf("expected key=value, got %q", secrets["key"])
	}
}

func TestKV2Client_DoesNotDoubleRewrite(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/data/myapp": {"token": "abc123"},
	})
	client := NewKV2Client(mock, "secret")

	// Path already contains /data/ — should not be rewritten again.
	secrets, err := client.ReadSecrets(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["token"] != "abc123" {
		t.Errorf("expected token=abc123, got %q", secrets["token"])
	}
}

func TestKV2Client_DefaultMount(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/data/cfg": {"db": "postgres"},
	})
	// No mounts provided — should default to "secret".
	client := NewKV2Client(mock)

	secrets, err := client.ReadSecrets(context.Background(), "secret/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["db"] != "postgres" {
		t.Errorf("expected db=postgres, got %q", secrets["db"])
	}
}

func TestKV2Client_NonMatchingMount(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"kv/myapp": {"host": "localhost"},
	})
	// Mount is "secret" but path uses "kv" — no rewrite should occur.
	client := NewKV2Client(mock, "secret")

	secrets, err := client.ReadSecrets(context.Background(), "kv/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", secrets["host"])
	}
}

func TestKV2Client_MultipleMounts(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"kvv2/data/svc": {"port": "8080"},
	})
	client := NewKV2Client(mock, "secret", "kvv2")

	secrets, err := client.ReadSecrets(context.Background(), "kvv2/svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", secrets["port"])
	}
}

func TestKV2Client_PathNotFound(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/data/myapp": {"key": "value"},
	})
	client := NewKV2Client(mock, "secret")

	// Reading a path that does not exist should return an error.
	_, err := client.ReadSecrets(context.Background(), "secret/nonexistent")
	if err == nil {
		t.Fatal("expected an error for missing path, got nil")
	}
}
