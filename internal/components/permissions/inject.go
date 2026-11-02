package permissions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zabadev/agent-ai/internal/agents"
	"github.com/zabadev/agent-ai/internal/components/filemerge"
	"github.com/zabadev/agent-ai/internal/model"
)

type InjectionResult struct {
	Changed bool
	Files   []string
}

// Permission overlays intentionally keep execution approval enabled by default.
// Tool-level secret denials add defense in depth; they do not sandbox shell access.
var claudeCodeOverlayJSON = []byte(`{
  "permissions": {
    "defaultMode": "default",
    "deny": [
      "Read(.env)",
      "Read(.env.*)",
      "Read(**/.env)",
      "Read(**/.env.*)",
      "Read(**/secrets/**)",
      "Read(**/credentials.json)",
      "Edit(.env)",
      "Edit(.env.*)",
      "Edit(**/.env)",
      "Edit(**/.env.*)",
      "Edit(**/secrets/**)",
      "Edit(**/credentials.json)"
    ]
  }
}
`)

// OpenCode asks before every shell command. Secret-path denials are defense in
// depth only: users should still use OS-level sandboxing for untrusted code.
var openCodeOverlayJSON = []byte(`{
  "permission": {
    "bash": {
      "*": "ask"
    },
    "read": {
      "*": "allow",
      "*.env": "deny",
      "*.env.*": "deny",
      "**/.env": "deny",
      "**/.env.*": "deny",
      "**/.envrc": "deny",
      "**/.npmrc": "deny",
      "**/secrets/**": "deny",
      "**/credentials.json": "deny",
      "**/*private*key*": "deny"
    },
    "edit": {
      "*.env": "deny",
      "*.env.*": "deny",
      "**/.env": "deny",
      "**/.env.*": "deny",
      "**/.envrc": "deny",
      "**/.npmrc": "deny",
      "**/secrets/**": "deny",
      "**/credentials.json": "deny",
      "**/*private*key*": "deny"
    }
  }
}
`)

var geminiCLIOverlayJSON = []byte(`{
  "general": {
    "defaultApprovalMode": "default"
  }
}
`)

var vscodeCopilotOverlayJSON = []byte(`{
  "chat.tools.autoApprove": false
}
`)

// agentOverlay returns the correct permission overlay for the given agent,
// or nil if the agent does not support permission injection via settings.json.
func agentOverlay(id model.AgentID) []byte {
	switch id {
	case model.AgentClaudeCode:
		return claudeCodeOverlayJSON
	case model.AgentOpenCode:
		return openCodeOverlayJSON
	case model.AgentGeminiCLI:
		return geminiCLIOverlayJSON
	case model.AgentVSCodeCopilot:
		return vscodeCopilotOverlayJSON
	case model.AgentCursor:
		// Cursor manages permissions via cli-config.json, not settings.json.
		return nil
	case model.AgentCodex:
		// Codex has no known settings.json path; permissions are skipped.
		return nil
	default:
		return nil
	}
}

func Inject(homeDir string, adapter agents.Adapter) (InjectionResult, error) {
	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		return InjectionResult{}, nil
	}

	overlay := agentOverlay(adapter.Agent())
	if overlay == nil {
		return InjectionResult{}, nil
	}

	writeResult, err := mergeJSONFile(settingsPath, overlay, adapter.Agent())
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
}

func mergeJSONFile(path string, overlay []byte, agent model.AgentID) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, err
	}
	if agent == model.AgentClaudeCode {
		merged, err = preserveClaudeDenyRules(baseJSON, merged)
		if err != nil {
			return filemerge.WriteResult{}, err
		}
	}

	return filemerge.WriteFileAtomic(path, merged, 0o644)
}

// preserveClaudeDenyRules unions our managed secret-path denials with user
// rules; replacing this array would silently remove custom protections.
func preserveClaudeDenyRules(existing, merged []byte) ([]byte, error) {
	if len(existing) == 0 {
		return merged, nil
	}
	var existingRoot, mergedRoot map[string]any
	if err := json.Unmarshal(existing, &existingRoot); err != nil {
		return nil, fmt.Errorf("parse existing Claude settings: %w", err)
	}
	if err := json.Unmarshal(merged, &mergedRoot); err != nil {
		return nil, fmt.Errorf("parse merged Claude settings: %w", err)
	}
	existingPermissions, _ := existingRoot["permissions"].(map[string]any)
	mergedPermissions, _ := mergedRoot["permissions"].(map[string]any)
	if existingPermissions == nil || mergedPermissions == nil {
		return merged, nil
	}
	oldRules, _ := existingPermissions["deny"].([]any)
	managedRules, _ := mergedPermissions["deny"].([]any)
	seen := make(map[string]struct{}, len(oldRules)+len(managedRules))
	rules := make([]any, 0, len(oldRules)+len(managedRules))
	for _, rule := range append(oldRules, managedRules...) {
		text, ok := rule.(string)
		if !ok {
			continue
		}
		if _, exists := seen[text]; exists {
			continue
		}
		seen[text] = struct{}{}
		rules = append(rules, text)
	}
	mergedPermissions["deny"] = rules
	result, err := json.MarshalIndent(mergedRoot, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal merged Claude settings: %w", err)
	}
	return append(result, '\n'), nil
}

var osReadFile = func(path string) ([]byte, error) {
	// #nosec G304 -- path is derived from user home directory, not external input
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}

	return content, nil
}
