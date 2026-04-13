package vault_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestTimeoutClient_SucceedsWithinTimeout(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/data": {"key": "value"},
	})
	c := vault.NewTimeoutClient(mock, 500*time.Millisecond)

	data, err := c.ReadSecrets(context.Background(), "secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("expected value %q, got %q", "value", data["key"])
	}
}

func TestTimeoutClient_ExceedsTimeout(t *testing.T) {
	slow := &slowReader{delay: 200 * time.Millisecond}
	c := vault.NewTimeoutClient(slow, 50*time.Millisecond)

	_, err := c.ReadSecrets(context.Background(), "secret/slow")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestTimeoutClient_ZeroTimeoutDisabled(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/data": {"k": "v"},
	})
	c := vault.NewTimeoutClient(mock, 0)

	data, err := c.ReadSecrets(context.Background(), "secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["k"] != "v" {
		t.Errorf("expected %q, got %q", "v", data["k"])
	}
}

func TestTimeoutClient_PropagatesInnerError(t *testing.T) {
	sentinel := errors.New("inner failure")
	mock := vault.NewMockClient(nil)
	mock.SetError("secret/fail", sentinel)
	c := vault.NewTimeoutClient(mock, 500*time.Millisecond)

	_, err := c.ReadSecrets(context.Background(), "secret/fail")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTimeoutClient_ImplementsSecretReader(t *testing.T) {
	var _ vault.SecretReader = vault.NewTimeoutClient(vault.NewMockClient(nil), time.Second)
}

// slowReader simulates a backend that takes a configurable time to respond.
type slowReader struct {
	delay time.Duration
}

func (s *slowReader) ReadSecrets(ctx context.Context, _ string) (map[string]string, error) {
	select {
	case <-time.After(s.delay):
		return map[string]string{"slow": "result"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
