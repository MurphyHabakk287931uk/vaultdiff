package diff

import (
	"context"
	"fmt"

	"github.com/user/vaultdiff/internal/redact"
	"github.com/user/vaultdiff/internal/vault"
)

// RunOptions configures a diff run between two Vault paths.
type RunOptions struct {
	SrcPath    string
	DstPath    string
	RedactMode redact.Mode
	ShowUnchanged bool
}

// Runner orchestrates reading secrets and producing a diff.
type Runner struct {
	client vault.SecretReader
}

// NewRunner creates a Runner with the given SecretReader.
func NewRunner(client vault.SecretReader) *Runner {
	return &Runner{client: client}
}

// Run reads secrets from both paths, compares them, and returns formatted output.
func (r *Runner) Run(ctx context.Context, opts RunOptions) (string, error) {
	src, err := r.client.ReadSecrets(ctx, opts.SrcPath)
	if err != nil {
		return "", fmt.Errorf("reading source path %q: %w", opts.SrcPath, err)
	}

	dst, err := r.client.ReadSecrets(ctx, opts.DstPath)
	if err != nil {
		return "", fmt.Errorf("reading destination path %q: %w", opts.DstPath, err)
	}

	result := Compare(src, dst)

	if !opts.ShowUnchanged {
		result = filterUnchanged(result)
	}

	return Format(result, opts.RedactMode), nil
}

func filterUnchanged(r *Result) *Result {
	var filtered []Entry
	for _, e := range r.Entries {
		if e.Change != Unchanged {
			filtered = append(filtered, e)
		}
	}
	return &Result{Entries: filtered}
}
