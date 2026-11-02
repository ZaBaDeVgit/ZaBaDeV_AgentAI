package permissions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/zabadev/agent-ai/internal/agents/claude"
)

func TestInjectClaudePreservesExistingDenyRules(t *testing.T) {
	home := t.TempDir()
	settingsPath := claude.NewAdapter().SettingsPath(home)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("create settings directory: %v", err)
	}
	before := `{"permissions":{"deny":["Bash(rm -rf /)","Read(.private-file)"]}}`
	if err := os.WriteFile(settingsPath, []byte(before), 0o600); err != nil {
		t.Fatalf("write existing settings: %v", err)
	}

	if _, err := Inject(home, claude.NewAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read merged settings: %v", err)
	}
	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("parse merged settings: %v", err)
	}
	permissionsNode := settings["permissions"].(map[string]any)
	deny := permissionsNode["deny"].([]any)
	got := make(map[string]bool, len(deny))
	for _, value := range deny {
		got[value.(string)] = true
	}
	for _, rule := range []string{"Bash(rm -rf /)", "Read(.private-file)", "Read(.env)", "Edit(**/secrets/**)"} {
		if !got[rule] {
			t.Errorf("deny rules lost %q; got %#v", rule, deny)
		}
	}
}
