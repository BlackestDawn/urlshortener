package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_AllEnvVarsSet_UsesProvidedValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("LISTEN_PORT", "9090")
	t.Setenv("SHORTEN_DOMAIN", "short.example.com")
	t.Setenv("URLSHORTENER_ENV", "test")

	cfg := NewConfig()

	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DBUrl)
	assert.Equal(t, ":9090", cfg.Port)
	assert.Equal(t, "short.example.com", cfg.Domain)
	assert.Equal(t, "test", cfg.Env)
}

func TestNewConfig_OptionalVarsUnset_FallsBackToDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("LISTEN_PORT", "")
	t.Setenv("SHORTEN_DOMAIN", "")
	t.Setenv("URLSHORTENER_ENV", "")

	cfg := NewConfig()

	assert.Equal(t, ":"+defaultListenPort, cfg.Port)
	assert.Equal(t, defaultDomain, cfg.Domain)
	assert.Equal(t, defaultAppEnv, cfg.Env)
}

func TestAddCloser_Cleanup_CallsAllClosers(t *testing.T) {
	cfg := &Config{}
	calledA := false
	calledB := false

	cfg.AddCloser(func() error {
		calledA = true
		return nil
	})
	cfg.AddCloser(func() error {
		calledB = true
		return errors.New("boom")
	})

	assert.NotPanics(t, cfg.Cleanup)
	assert.True(t, calledA)
	assert.True(t, calledB)
}

func TestFindEnvFile_FileExistsInParentDir_ReturnsPath(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b")
	assert.NoError(t, os.MkdirAll(sub, 0o755))

	envPath := filepath.Join(root, ".env.test")
	assert.NoError(t, os.WriteFile(envPath, []byte(""), 0o644))

	wd, err := os.Getwd()
	assert.NoError(t, err)
	t.Chdir(sub)
	defer t.Chdir(wd)

	got := findEnvFile("test")
	assert.Equal(t, envPath, got)
}

func TestFindEnvFile_NoMatchingFile_ReturnsEmptyString(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	assert.NoError(t, err)
	t.Chdir(dir)
	defer t.Chdir(wd)

	got := findEnvFile("nonexistent")
	assert.Equal(t, "", got)
}
