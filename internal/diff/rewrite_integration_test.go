package diff_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/redact"
	"github.com/your-org/vaultdiff/internal/vault"
)

func TestRunner_RewriteClientRemapsSourcePath(t *testing.T) {
	srcData := map[string]map[string]string{
		"prod/db": {"password": "hunter2"},
	}
	dstData := map[string]map[string]string{
		"prod/db": {"password": "hunter3"},
	}

	srcBase := vault.NewMockClient(srcData)
	src := vault.NewRewriteClient(srcBase, []vault.RewriteRule{
		{From: "staging/", To: "prod/"},
	})
	dst := vault.NewMockClient(dstData)

	r := diff.NewRunner(src, dst, redact.ModeRedact, false)
	results, err := r.Run("staging/db", "prod/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one diff result")
	}
	formatted := diff.Format(results)
	if !strings.Contains(formatted, "password") {
		t.Errorf("expected formatted diff to mention password key, got:\n%s", formatted)
	}
}
