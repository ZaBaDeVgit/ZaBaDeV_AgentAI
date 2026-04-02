package vscode

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/zabadev/agent-ai/internal/model"
	"github.com/zabadev/agent-ai/internal/system"
)

func TestStrategies(t *testing.T) {
	a := NewAdapter()

	if got := a.SystemPromptStrategy(); got != model.StrategyInstructionsFile {
		t.Fatalf("SystemPromptStrategy() = %v, want %v", got, model.StrategyInstructionsFile)
	}

	if got := a.MCPStrategy(); got != model.StrategyMCPConfigFile {
		t.Fatalf("MCPStrategy() = %v, want %v", got, model.StrategyMCPConfigFile)
	}
}

func TestSystemPromptFileUsesInstructionsExtension(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	path := a.SystemPromptFile(home)
	if filepath.Ext(path) != ".md" {
		t.Fatalf("SystemPromptFile() should end with .md: %q", path)
	}

	if filepath.Base(path) != "gentle-ai.instructions.md" {
		t.Fatalf("SystemPromptFile() = %q, want filename gentle-ai.instructions.md", path)
	}
}

func TestSettingsPathUsesVSCodeUserProfile(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	switch runtime.GOOS {
	case "darwin":
		path := a.SettingsPath(home)
		want := filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	case "windows":
		appData := filepath.Join(home, "AppData", "Roaming")
		t.Setenv("APPDATA", appData)
		path := a.SettingsPath(home)
		want := filepath.Join(appData, "Code", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	default:
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg"))
		path := a.SettingsPath(home)
		want := filepath.Join(home, "xdg", "Code", "User", "settings.json")
		if path != want {
			t.Fatalf("SettingsPath() = %q, want %q", path, want)
		}
	}
}

func TestMCPConfigPathUsesVSCodeUserProfile(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	switch runtime.GOOS {
	case "darwin":
		path := a.MCPConfigPath(home, "context7")
		want := filepath.Join(home, "Library", "Application Support", "Code", "User", "mcp.json")
		if path != want {
			t.Fatalf("MCPConfigPath() = %q, want %q", path, want)
		}
	case "windows":
		appData := filepath.Join(home, "AppData", "Roaming")
		t.Setenv("APPDATA", appData)
		path := a.MCPConfigPath(home, "context7")
		want := filepath.Join(appData, "Code", "User", "mcp.json")
		if path != want {
			t.Fatalf("MCPConfigPath() = %q, want %q", path, want)
		}
	default:
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg"))
		path := a.MCPConfigPath(home, "context7")
		want := filepath.Join(home, "xdg", "Code", "User", "mcp.json")
		if path != want {
			t.Fatalf("MCPConfigPath() = %q, want %q", path, want)
		}
	}
}

func TestAdapterIdentity(t *testing.T) {
	tests := []struct {
		name      string
		wantAgent model.AgentID
		wantTier  model.SupportTier
	}{
		{
			name:      "vscode agent identity",
			wantAgent: model.AgentVSCodeCopilot,
			wantTier:  model.TierFull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			if got := a.Agent(); got != tt.wantAgent {
				t.Errorf("Agent() = %v, want %v", got, tt.wantAgent)
			}
			if got := a.Tier(); got != tt.wantTier {
				t.Errorf("Tier() = %v, want %v", got, tt.wantTier)
			}
		})
	}
}

func TestSystemPromptFileTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home directory",
			homeDir:  "/home/user",
			expected: "gentle-ai.instructions.md",
		},
		{
			name:     "path with special characters",
			homeDir:  "/home/user-name_123",
			expected: "gentle-ai.instructions.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			got := a.SystemPromptFile(tt.homeDir)
			// Just verify it ends with the expected filename
			if filepath.Base(got) != tt.expected {
				t.Errorf("SystemPromptFile() basename = %q, want %q", filepath.Base(got), tt.expected)
			}
		})
	}
}

func TestGlobalConfigDir(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.copilot",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.copilot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.GlobalConfigDir(tt.homeDir)
			if got != tt.expected {
				t.Errorf("GlobalConfigDir() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSkillsDir(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.copilot/skills",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.copilot/skills",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.SkillsDir(tt.homeDir)
			if got != tt.expected {
				t.Errorf("SkillsDir() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSystemPromptDir(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	// On Linux (default), uses XDG_CONFIG_HOME
	if runtime.GOOS == "linux" {
		t.Setenv("XDG_CONFIG_HOME", "/home/user/.config")
		got := a.SystemPromptDir(home)
		expected := filepath.Join("/home/user/.config", "Code", "User", "prompts")
		if got != expected {
			t.Errorf("SystemPromptDir() = %q, want %q", got, expected)
		}
	}
}

func TestCapabilities(t *testing.T) {
	tests := []struct {
		name              string
		wantOutputStyles  bool
		wantSlashCommands bool
		wantSkills        bool
		wantSystemPrompt  bool
		wantMCP           bool
		wantAutoInstall   bool
	}{
		{
			name:              "vscode capabilities",
			wantOutputStyles:  false,
			wantSlashCommands: false,
			wantSkills:        true,
			wantSystemPrompt:  true,
			wantMCP:           true,
			wantAutoInstall:   false, // Desktop app, cannot auto-install
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			if got := a.SupportsOutputStyles(); got != tt.wantOutputStyles {
				t.Errorf("SupportsOutputStyles() = %v, want %v", got, tt.wantOutputStyles)
			}
			if got := a.SupportsSlashCommands(); got != tt.wantSlashCommands {
				t.Errorf("SupportsSlashCommands() = %v, want %v", got, tt.wantSlashCommands)
			}
			if got := a.SupportsSkills(); got != tt.wantSkills {
				t.Errorf("SupportsSkills() = %v, want %v", got, tt.wantSkills)
			}
			if got := a.SupportsSystemPrompt(); got != tt.wantSystemPrompt {
				t.Errorf("SupportsSystemPrompt() = %v, want %v", got, tt.wantSystemPrompt)
			}
			if got := a.SupportsMCP(); got != tt.wantMCP {
				t.Errorf("SupportsMCP() = %v, want %v", got, tt.wantMCP)
			}
			if got := a.SupportsAutoInstall(); got != tt.wantAutoInstall {
				t.Errorf("SupportsAutoInstall() = %v, want %v", got, tt.wantAutoInstall)
			}
		})
	}
}

func TestCommandsDir(t *testing.T) {
	a := NewAdapter()
	// VS Code doesn't support slash commands
	if got := a.CommandsDir("/home/user"); got != "" {
		t.Errorf("CommandsDir() = %q, want empty string", got)
	}
}

func TestOutputStyleDir(t *testing.T) {
	a := NewAdapter()
	// VS Code doesn't support output styles
	if got := a.OutputStyleDir("/home/user"); got != "" {
		t.Errorf("OutputStyleDir() = %q, want empty string", got)
	}
}

func TestInstallCommandNotInstallable(t *testing.T) {
	a := NewAdapter()

	_, err := a.InstallCommand(system.PlatformProfile{})
	if err == nil {
		t.Fatalf("InstallCommand() should return error for desktop app")
	}

	var installErr AgentNotInstallableError
	if !errors.As(err, &installErr) {
		t.Fatalf("InstallCommand() error = %v, want AgentNotInstallableError", err)
	}
}
