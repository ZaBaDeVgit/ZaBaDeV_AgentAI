package components_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zabadev/agent-ai/internal/agents/opencode"
	"github.com/zabadev/agent-ai/internal/components/engram"
	"github.com/zabadev/agent-ai/internal/components/mcp"
	"github.com/zabadev/agent-ai/internal/components/permissions"
	"github.com/zabadev/agent-ai/internal/components/persona"
	"github.com/zabadev/agent-ai/internal/model"
)

// This integration test checks the installed OpenCode config and instruction
// file together: a secure approval policy is useful only if it is actually
// emitted alongside the configured agent that loads those instructions.
func TestOpenCodeInstallWiresInstructionsAndSafePermissions(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	if _, err := persona.Inject(home, adapter, model.PersonaSeniorZaBaDeV); err != nil {
		t.Fatalf("inject persona: %v", err)
	}
	if _, err := permissions.Inject(home, adapter); err != nil {
		t.Fatalf("inject permissions: %v", err)
	}
	if _, err := engram.Inject(home, adapter); err != nil {
		t.Fatalf("inject Engram: %v", err)
	}
	if _, err := mcp.Inject(home, adapter); err != nil {
		t.Fatalf("inject Context7: %v", err)
	}

	settingsPath := adapter.SettingsPath(home)
	settingsBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read OpenCode settings: %v", err)
	}
	var settings map[string]any
	if err := json.Unmarshal(settingsBytes, &settings); err != nil {
		t.Fatalf("parse OpenCode settings: %v", err)
	}

	agents, ok := settings["agent"].(map[string]any)
	if !ok {
		t.Fatalf("agent definitions missing: %#v", settings["agent"])
	}
	senior, ok := agents["senior-zabadev"].(map[string]any)
	if !ok {
		t.Fatalf("senior-zabadev agent missing: %#v", agents["senior-zabadev"])
	}
	if senior["prompt"] != "{file:./AGENTS.md}" {
		t.Fatalf("agent prompt = %v, want AGENTS.md reference", senior["prompt"])
	}

	permission, ok := settings["permission"].(map[string]any)
	if !ok {
		t.Fatalf("permission policy missing: %#v", settings["permission"])
	}
	bash, ok := permission["bash"].(map[string]any)
	if !ok || bash["*"] != "ask" {
		t.Fatalf("shell wildcard must require approval, got %#v", permission["bash"])
	}

	read, ok := permission["read"].(map[string]any)
	if !ok {
		t.Fatalf("read policy missing: %#v", permission["read"])
	}
	for _, secretPattern := range []string{"**/.env", "**/.env.*", "**/secrets/**", "**/credentials.json"} {
		if read[secretPattern] != "deny" {
			t.Errorf("read policy for %q = %v, want deny", secretPattern, read[secretPattern])
		}
	}

	mcpServers, ok := settings["mcp"].(map[string]any)
	if !ok {
		t.Fatalf("MCP server configuration missing: %#v", settings["mcp"])
	}
	engramConfig, ok := mcpServers["engram"].(map[string]any)
	if !ok || engramConfig["enabled"] != true {
		t.Fatalf("Engram MCP server not enabled: %#v", mcpServers["engram"])
	}
	command, ok := engramConfig["command"].([]any)
	if !ok || len(command) < 3 || command[0] != "engram" || command[1] != "mcp" {
		t.Fatalf("unexpected Engram MCP command: %#v", engramConfig["command"])
	}
	context7Config, ok := mcpServers["context7"].(map[string]any)
	if !ok || context7Config["type"] != "remote" || context7Config["enabled"] != true {
		t.Fatalf("Context7 MCP server not configured: %#v", mcpServers["context7"])
	}

	promptPath := filepath.Join(adapter.SystemPromptDir(home), "AGENTS.md")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("read injected agent instructions: %v", err)
	}
	prompt := string(promptBytes)
	for _, expected := range []string{"Never reveal, copy or transmit secrets", "Run the most relevant checks", "Never say a test passed if it was not executed"} {
		if !strings.Contains(prompt, expected) {
			t.Errorf("injected instructions missing safety/verification rule %q", expected)
		}
	}
}
