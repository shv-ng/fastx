package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePyproject(t *testing.T) {
	tomlContent := `
	project="hello"

[tool.fastx]
version = "v0.1"
registry = "https://github.com/shv-ng/fastx"
skipped = ["auth","logging"]

[tool.fastx.components]
database = "v0.1"
cache = "v0.1"
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "pyproject.toml")
	err := os.WriteFile(tmpFile, []byte(tomlContent), 0644)
	require.NoError(t, err, "fail to write mock toml file")

	meta, err := ParsePyproject(tmpFile)
	require.NoError(t, err, "ParsePyproject() should not return any error")
	require.NotNil(t, meta, "meta should not be nil")

	assert.Equal(t, "v0.1", meta.Version)
	assert.Equal(t, "https://github.com/shv-ng/fastx", meta.Registry)
	assert.Equal(t, []string{"auth", "logging"}, meta.Skipped)
	assert.Equal(t, map[string]string{"database": "v0.1", "cache": "v0.1"}, meta.Components)

}

func TestWriteFastxMeta(t *testing.T) {
	tomlContent := `
	project="hello"

[tool.fastx]
version = "v0.1"
registry = "https://github.com/shv-ng/fastx"
skipped = ["auth","logging"]

[tool.fastx.components]
database = "v0.1"
cache = "v0.1"
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "pyproject.toml")
	err := os.WriteFile(tmpFile, []byte(tomlContent), 0644)
	require.NoError(t, err, "fail to write mock toml file")

	upsertMeta := &FastxMeta{
		Version:    "v0.2",
		Registry:   "https://github.com/shv-ng/fastx/new",
		Skipped:    []string{},
		Components: map[string]string{},
	}

	err = WriteFastxMeta(tmpFile, upsertMeta)
	require.NoError(t, err, "WriteFastxMeta() should not return any error")

	newMeta, err := ParsePyproject(tmpFile)
	require.NoError(t, err, "ParsePyproject() should not return any error")

	assert.Equal(t, upsertMeta, newMeta)
}
