package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "diff <src-path> <dst-path>" {
			found = true
			break
		}
	}
	assert.True(t, found, "diff subcommand should be registered")
}

func TestDiffCmd_RequiresTwoArgs(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"diff", "secret/only-one"})
	buf := &bytes.Buffer{}
	cmd.SetErr(buf)
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestDiffCmd_InvalidRedactMode(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"diff", "--redact=bogus", "secret/a", "secret/b"})
	buf := &bytes.Buffer{}
	cmd.SetErr(buf)
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestDiffCmd_InvalidOutputFormat(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"diff", "--output=yaml", "secret/a", "secret/b"})
	buf := &bytes.Buffer{}
	cmd.SetErr(buf)
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestRootCmd_DefaultFlags(t *testing.T) {
	f := rootCmd.PersistentFlags()

	v, err := f.GetString("redact")
	require.NoError(t, err)
	assert.Equal(t, "none", v)

	o, err := f.GetString("output")
	require.NoError(t, err)
	assert.Equal(t, "text", o)

	a, err := f.GetBool("all")
	require.NoError(t, err)
	assert.False(t, a)
}
