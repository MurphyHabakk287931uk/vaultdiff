package cmd

import (
	"fmt"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// buildSplitClient constructs a SplitClient when --split mode is requested.
// leftNS and rightNS are the namespace labels applied to each side's keys.
func buildSplitClient(
	left vault.SecretReader,
	leftNS string,
	right vault.SecretReader,
	rightNS string,
) (vault.SecretReader, error) {
	if leftNS == "" {
		return nil, fmt.Errorf("split: left namespace must not be empty")
	}
	if rightNS == "" {
		return nil, fmt.Errorf("split: right namespace must not be empty")
	}
	return vault.NewSplitClient(left, leftNS, right, rightNS), nil
}
