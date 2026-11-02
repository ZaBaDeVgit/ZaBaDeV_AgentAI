package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zabadev/agent-ai/internal/model"
)

func TestVerifyOpenCodeConfigConfirmsConfiguredFeatures(t *testing.T) {
	home := t.TempDir()
	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	pluginPath := filepath.Join(home, ".config", "opencode", "plugins", "background-agents.ts")
	instructionsPath := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	for _, path := range []string{pluginPath, instructionsPath} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create test directory: %v", err)
		}
		if err := os.WriteFile(path, []byte("test asset"), 0o600); err != nil {
			t.Fatalf("write test asset: %v", err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("create settings directory: %v", err)
	}
	config := `{
  "agent": {"senior-zabadev": {"prompt": "{file:./AGENTS.md}"}},
  "permission": {"bash": {"*": "ask"}},
  "mcp": {
    "engram": {"enabled": true},
    "context7": {"enabled": true}
  }
}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	components := []model.ComponentID{
		model.ComponentPersona,
		model.ComponentPermission,
		model.ComponentEngram,
		model.ComponentContext7,
		model.ComponentSDD,
	}
	if err := verifyOpenCodeConfig(configPath, home, components, model.PersonaSeniorZaBaDeV); err != nil {
		t.Fatalf("verifyOpenCodeConfig() error = %v", err)
	}
}

func TestVerifyOpenCodeConfigRejectsUnsafePermissionPolicy(t *testing.T) {
	home := t.TempDir()
	configPath := filepath.Join(home, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"permission":{"bash":{"*":"allow"}}}`), 0o600); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	err := verifyOpenCodeConfig(configPath, home, []model.ComponentID{model.ComponentPermission}, model.PersonaSeniorZaBaDeV)
	if err == nil {
		t.Fatal("verifyOpenCodeConfig() accepted shell wildcard allow")
	}
}

func TestVerifyOpenCodeConfigRejectsInvalidJSON(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "opencode.json")
	if err := os.WriteFile(configPath, []byte("{broken"), 0o600); err != nil {
		t.Fatalf("write settings: %v", err)
	}
	if err := verifyOpenCodeConfig(configPath, t.TempDir(), nil, model.PersonaCustom); err == nil {
		t.Fatal("verifyOpenCodeConfig() accepted malformed JSON")
	}
}
