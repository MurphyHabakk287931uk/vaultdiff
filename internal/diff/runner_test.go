package diff_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/user/vaultdiff/internal/diff"
	"github.com/user/vaultdiff/internal/redact"
	"github.com/user/vaultdiff/internal/vault"
)

func makeRunner(data map[string]map[string]string) *diff.Runner {
	mc := vault.NewMockClient(data)
	return diff.NewRunner(mc)
}

func TestRunner_Run_BasicDiff(t *testing.T) {
	data := map[string]map[string]string{
		"secret/src": {"key": "old"},
		"secret/dst": {"key": "new"},
	}
	r := makeRunner(data)
	out, err := r.Run(context.Background(), diff.RunOptions{
		SrcPath:    "secret/src",
		DstPath:    "secret/dst",
		RedactMode: redact.ModeNone,
		ShowUnchanged: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}

func TestRunner_Run_HidesUnchanged(t *testing.T) {
	data := map[string]map[string]string{
		"secret/a": {"same": "val", "diff": "v1"},
		"secret/b": {"same": "val", "diff": "v2"},
	}
	r := makeRunner(data)
	out, err := r.Run(context.Background(), diff.RunOptions{
		SrcPath:       "secret/a",
		DstPath:       "secret/b",
		RedactMode:    redact.ModeNone,
		ShowUnchanged: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "same") {
		t.Errorf("expected unchanged key to be hidden, got: %s", out)
	}
	if !strings.Contains(out, "diff") {
		t.Errorf("expected changed key to appear, got: %s", out)
	}
}

func TestRunner_Run_SrcError(t *testing.T) {
	mc := vault.NewMockClient(nil)
	mc.SetError("secret/src", errors.New("permission denied"))
	r := diff.NewRunner(mc)
	_, err := r.Run(context.Background(), diff.RunOptions{
		SrcPath: "secret/src",
		DstPath: "secret/dst",
	})
	if err == nil {
		t.Fatal("expected error from src read")
	}
	if !strings.Contains(err.Error(), "source path") {
		t.Errorf("expected source path in error, got: %v", err)
	}
}

func TestRunner_Run_MaskMode(t *testing.T) {
	data := map[string]map[string]string{
		"secret/x": {"token": "supersecret"},
		"secret/y": {"token": "othersecret"},
	}
	r := makeRunner(data)
	out, err := r.Run(context.Background(), diff.RunOptions{
		SrcPath:       "secret/x",
		DstPath:       "secret/y",
		RedactMode:    redact.ModeMask,
		ShowUnchanged: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "supersecret") || strings.Contains(out, "othersecret") {
		t.Errorf("expected masked output, got: %s", out)
	}
}
