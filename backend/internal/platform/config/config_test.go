package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFptAPIKeyFromSecretFile(t *testing.T) {
	t.Setenv("FPT_AI_API_KEY", "")
	secretPath := filepath.Join(t.TempDir(), "fpt-key")
	if err := os.WriteFile(secretPath, []byte("secret-from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FPT_AI_API_KEY_FILE", secretPath)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FptApiKey != "secret-from-file" {
		t.Fatalf("FptApiKey = %q, want secret-from-file", cfg.FptApiKey)
	}
}

func TestLoadExplicitFptAPIKeyOverridesSecretFile(t *testing.T) {
	t.Setenv("FPT_AI_API_KEY", "explicit-key")
	t.Setenv("FPT_AI_API_KEY_FILE", "/missing/secret")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FptApiKey != "explicit-key" {
		t.Fatalf("FptApiKey = %q, want explicit-key", cfg.FptApiKey)
	}
}

func TestLoadReturnsErrorForMissingConfiguredSecretFile(t *testing.T) {
	t.Setenv("FPT_AI_API_KEY", "")
	t.Setenv("FPT_AI_API_KEY_FILE", "/missing/secret")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want missing secret file error")
	}
}
