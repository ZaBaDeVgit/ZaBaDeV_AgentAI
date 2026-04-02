package persona

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zabadev/agent-ai/internal/agents"
	"github.com/zabadev/agent-ai/internal/assets"
	"github.com/zabadev/agent-ai/internal/components/filemerge"
	"github.com/zabadev/agent-ai/internal/model"
)

type InjectionResult struct {
	Changed bool
	Files   []string
}

const neutralPersonaContent = "Be helpful, direct, and technically precise. Focus on accuracy and clarity.\n"

// outputStyleOverlayJSON is the settings.json overlay to enable the ZaBaDeV output style.
var outputStyleOverlayJSON = []byte("{\n  \"outputStyle\": \"ZaBaDeV\"\n}\n")

const openCodeSeniorAgentID = "senior-zabadev"

func openCodeAgentOverlayJSON(_ model.PersonaID) []byte {
	return []byte(fmt.Sprintf("{\n  \"agent\": {\n    %q: {\n      \"mode\": \"primary\",\n      \"description\": \"Senior ZaBaDeV - single public agent with senior execution, memory, and SDD orchestration\",\n      \"prompt\": \"{file:./AGENTS.md}\",\n      \"tools\": {\n        \"read\": true,\n        \"write\": true,\n        \"edit\": true,\n        \"bash\": true,\n        \"delegate\": true,\n        \"delegation_read\": true,\n        \"delegation_list\": true\n      }\n    }\n  }\n}\n", openCodeSeniorAgentID))
}

func Inject(homeDir string, adapter agents.Adapter, persona model.PersonaID) (InjectionResult, error) {
	if !adapter.SupportsSystemPrompt() {
		return InjectionResult{}, nil
	}

	// Custom persona does nothing — user keeps their own config.
	if persona == model.PersonaCustom {
		return InjectionResult{}, nil
	}

	files := make([]string, 0, 2)
	changed := false

	content := personaContent(adapter.Agent(), persona)
	if content == "" {
		return InjectionResult{}, nil
	}

	// 1. Inject persona content based on system prompt strategy.
	switch adapter.SystemPromptStrategy() {
	case model.StrategyMarkdownSections:
		promptPath := adapter.SystemPromptFile(homeDir)
		existing, err := readFileOrEmpty(promptPath)
		if err != nil {
			return InjectionResult{}, err
		}

		// Auto-heal: strip any legacy free-text Gentleman persona block that was
		// written before the marker-based injection system existed. This prevents
		// duplicate persona content when users re-run the installer after an old
		// install placed the persona as raw text above the <!-- gentle-ai: --> markers.
		healed := filemerge.StripLegacyPersonaBlock(existing)

		updated := filemerge.InjectMarkdownSection(healed, "persona", content)

		writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || writeResult.Changed
		files = append(files, promptPath)

	case model.StrategyFileReplace:
		promptPath := adapter.SystemPromptFile(homeDir)
		writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(content), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || writeResult.Changed
		files = append(files, promptPath)

	case model.StrategyInstructionsFile:
		promptPath := adapter.SystemPromptFile(homeDir)

		// Auto-heal: remove any stale Gentleman persona content left at the
		// old VSCode path (~/.github/copilot-instructions.md) that was written
		// by an older installer version.  VS Code still reads that path for
		// global instructions, so the two files would conflict.
		if cleaned, cleanErr := cleanLegacyVSCodePersona(homeDir); cleanErr == nil && cleaned {
			changed = true
		}

		// Write the new instructions file (with YAML frontmatter) to the current path.
		// WriteFileAtomic compares bytes, so it is naturally idempotent: it rewrites
		// whenever the on-disk content differs from instructionsContent, which covers
		// the case where an older install wrote persona content without frontmatter.
		instructionsContent := wrapInstructionsFile(content)
		writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(instructionsContent), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || writeResult.Changed
		files = append(files, promptPath)

	case model.StrategyAppendToFile:
		promptPath := adapter.SystemPromptFile(homeDir)
		writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(content), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || writeResult.Changed
		files = append(files, promptPath)
	}

	// 2. OpenCode agent definitions — Tab-switchable agents in opencode.json.
	if adapter.Agent() == model.AgentOpenCode && persona != model.PersonaCustom {
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath != "" {
			agentResult, err := mergeJSONFile(settingsPath, openCodeAgentOverlayJSON(persona))
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || agentResult.Changed
			files = append(files, settingsPath)
		}
	}

	// 3. Gentleman-only: write output style + merge into settings (if agent supports it).
	if persona == model.PersonaGentleman && adapter.SupportsOutputStyles() {
		outputStyleDir := adapter.OutputStyleDir(homeDir)
		if outputStyleDir != "" {
			outputStylePath := outputStyleDir + "/gentleman.md"
			outputStyleContent := assets.MustRead("claude/output-style-gentleman.md")

			styleResult, err := filemerge.WriteFileAtomic(outputStylePath, []byte(outputStyleContent), 0o644)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || styleResult.Changed
			files = append(files, outputStylePath)
		}

		// Merge "outputStyle": "ZaBaDeV" into settings.
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath != "" {
			settingsResult, err := mergeJSONFile(settingsPath, outputStyleOverlayJSON)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || settingsResult.Changed
			files = append(files, settingsPath)
		}
	}

	return InjectionResult{Changed: changed, Files: files}, nil
}

func personaContent(agent model.AgentID, persona model.PersonaID) string {
	switch persona {
	case model.PersonaNeutral:
		return neutralPersonaContent
	case model.PersonaSeniorZaBaDeV:
		switch agent {
		case model.AgentOpenCode:
			return assets.MustRead("opencode/persona-senior-zabadev.md")
		default:
			return assets.MustRead("generic/persona-senior-zabadev.md")
		}
	case model.PersonaCustom:
		return ""
	default:
		// Gentleman persona — try agent-specific asset, then generic fallback.
		switch agent {
		case model.AgentClaudeCode:
			return assets.MustRead("claude/persona-gentleman.md")
		case model.AgentOpenCode:
			return assets.MustRead("opencode/persona-gentleman.md")
		default:
			// Generic persona includes ZaBaDeV personality + skills table + SDD orchestrator.
			// Used by Gemini CLI, Cursor, VS Code Copilot, and any future agents.
			return assets.MustRead("generic/persona-gentleman.md")
		}
	}
}

func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	return filemerge.WriteFileAtomic(path, merged, 0o644)
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

func readFileOrEmpty(path string) (string, error) {
	// #nosec G304 -- path is derived from user home directory, not external input
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(data), nil
}

func wrapInstructionsFile(content string) string {
	frontmatter := "---\n" +
		"name: ZaBaDeV Persona\n" +
		"description: ZaBaDeV persona with SDD orchestration and Engram protocol\n" +
		"applyTo: \"**\"\n" +
		"---\n\n"

	return frontmatter + content
}

// isLegacyUnwrappedPersona reports whether content looks like a Gentleman persona
// file that was written without YAML frontmatter by an older installer version.
// It returns true when the content carries known persona fingerprints but does NOT
// start with the YAML front-matter block ("---\n").
func isLegacyUnwrappedPersona(content string) bool {
	if strings.HasPrefix(content, "---\n") {
		// Already has YAML frontmatter — not a legacy file.
		return false
	}
	// Must contain at least one characteristic persona fingerprint.
	personaFingerprints := []string{
		"## Personality",
		"Senior Architect",
	}
	for _, fp := range personaFingerprints {
		if strings.Contains(content, fp) {
			return true
		}
	}
	return false
}

// legacyVSCodePersonaPaths returns the old VS Code persona file paths that may
// contain stale Gentleman persona content from older installer versions.
// These paths are no longer written by the current installer but may still
// be read by VS Code, causing conflicting instructions.
func legacyVSCodePersonaPaths(homeDir string) []string {
	return []string{
		// v1 path: wrote raw persona to ~/.github/copilot-instructions.md
		filepath.Join(homeDir, ".github", "copilot-instructions.md"),
	}
}

// cleanLegacyVSCodePersona removes Gentleman persona content from any old VS Code
// persona file paths that are no longer written by the current installer.
// Only files that contain clear Gentleman persona fingerprints are removed —
// files with user-written content are left untouched.
// Returns true if at least one file was cleaned.
func cleanLegacyVSCodePersona(homeDir string) (bool, error) {
	cleaned := false
	for _, oldPath := range legacyVSCodePersonaPaths(homeDir) {
		// #nosec G304 -- path is derived from user home directory, not external input
		data, err := os.ReadFile(oldPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return cleaned, fmt.Errorf("read legacy vscode persona %q: %w", oldPath, err)
		}

		if !isLegacyUnwrappedPersona(string(data)) {
			// File exists but doesn't look like a Gentleman persona — leave it alone.
			continue
		}

		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return cleaned, fmt.Errorf("remove legacy vscode persona %q: %w", oldPath, err)
		}
		cleaned = true
	}
	return cleaned, nil
}
