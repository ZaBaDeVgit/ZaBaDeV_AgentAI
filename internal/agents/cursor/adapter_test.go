package cursor

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/zabadev/agent-ai/internal/model"
	"github.com/zabadev/agent-ai/internal/system"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name            string
		stat            statResult
		wantInstalled   bool
		wantConfigPath  string
		wantConfigFound bool
		wantErr         bool
	}{
		{
			name:            "config directory found",
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantConfigPath:  filepath.Join("/tmp/home", ".cursor"),
			wantConfigFound: true,
		},
		{
			name:            "config missing",
			stat:            statResult{err: os.ErrNotExist},
			wantInstalled:   false,
			wantConfigPath:  filepath.Join("/tmp/home", ".cursor"),
			wantConfigFound: false,
		},
		{
			name:    "stat error bubbles up",
			stat:    statResult{err: errors.New("permission denied")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adapter{
				statPath: func(string) statResult {
					return tt.stat
				},
			}

			installed, _, configPath, configFound, err := a.Detect(context.Background(), "/tmp/home")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Detect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if installed != tt.wantInstalled {
				t.Fatalf("Detect() installed = %v, want %v", installed, tt.wantInstalled)
			}

			if configPath != tt.wantConfigPath {
				t.Fatalf("Detect() configPath = %q, want %q", configPath, tt.wantConfigPath)
			}

			if configFound != tt.wantConfigFound {
				t.Fatalf("Detect() configFound = %v, want %v", configFound, tt.wantConfigFound)
			}
		})
	}
}

func TestConfigPathsCrossPlatform(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	if got := a.GlobalConfigDir(home); got != filepath.Join(home, ".cursor") {
		t.Fatalf("GlobalConfigDir() = %q, want %q", got, filepath.Join(home, ".cursor"))
	}

	if got := a.SkillsDir(home); got != filepath.Join(home, ".cursor", "skills") {
		t.Fatalf("SkillsDir() = %q, want %q", got, filepath.Join(home, ".cursor", "skills"))
	}

	if got := a.MCPConfigPath(home, "ctx7"); got != filepath.Join(home, ".cursor", "mcp.json") {
		t.Fatalf("MCPConfigPath() = %q, want %q", got, filepath.Join(home, ".cursor", "mcp.json"))
	}

	if got := a.SystemPromptFile(home); got != filepath.Join(home, ".cursor", "rules", "gentle-ai.mdc") {
		t.Fatalf("SystemPromptFile() = %q, want %q", got, filepath.Join(home, ".cursor", "rules", "gentle-ai.mdc"))
	}
}

func TestStrategies(t *testing.T) {
	a := NewAdapter()

	if got := a.SystemPromptStrategy(); got != model.StrategyFileReplace {
		t.Fatalf("SystemPromptStrategy() = %v, want %v", got, model.StrategyFileReplace)
	}

	if got := a.MCPStrategy(); got != model.StrategyMCPConfigFile {
		t.Fatalf("MCPStrategy() = %v, want %v", got, model.StrategyMCPConfigFile)
	}
}

func TestDesktopAppNotAutoInstallable(t *testing.T) {
	a := NewAdapter()

	if a.SupportsAutoInstall() {
		t.Fatalf("Cursor should not support auto-install (desktop app)")
	}

	_, err := a.InstallCommand(system.PlatformProfile{})
	if err == nil {
		t.Fatalf("InstallCommand() should return error for desktop app")
	}
}

func TestAdapterIdentity(t *testing.T) {
	tests := []struct {
		name      string
		wantAgent model.AgentID
		wantTier  model.SupportTier
	}{
		{
			name:      "cursor agent identity",
			wantAgent: model.AgentCursor,
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

func TestSystemPromptPaths(t *testing.T) {
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home directory",
			homeDir:  "/home/user",
			expected: "/home/user/.cursor/rules/gentle-ai.mdc",
		},
		{
			name:     "path with special characters",
			homeDir:  "/home/user-name_123",
			expected: "/home/user-name_123/.cursor/rules/gentle-ai.mdc",
		},
		{
			name:     "empty home dir returns relative path",
			homeDir:  "",
			expected: ".cursor/rules/gentle-ai.mdc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			got := a.SystemPromptFile(tt.homeDir)
			if got != tt.expected {
				t.Errorf("SystemPromptFile() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSystemPromptDir(t *testing.T) {
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.cursor/rules",
		},
		{
			name:     "empty home returns relative path",
			homeDir:  "",
			expected: ".cursor/rules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			got := a.SystemPromptDir(tt.homeDir)
			if got != tt.expected {
				t.Errorf("SystemPromptDir() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGlobalConfigDirTableDriven(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.cursor",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.cursor",
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

func TestSkillsDirTableDriven(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.cursor/skills",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.cursor/skills",
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

func TestSettingsPathTableDriven(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.cursor/settings.json",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.cursor/settings.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.SettingsPath(tt.homeDir)
			if got != tt.expected {
				t.Errorf("SettingsPath() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestMCPConfigPathTableDriven(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name       string
		homeDir    string
		serverName string
		expected   string
	}{
		{
			name:       "context7 server",
			homeDir:    "/home/user",
			serverName: "context7",
			expected:   "/home/user/.cursor/mcp.json",
		},
		{
			name:       "filesystem server",
			homeDir:    "/home/user",
			serverName: "filesystem",
			expected:   "/home/user/.cursor/mcp.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.MCPConfigPath(tt.homeDir, tt.serverName)
			if got != tt.expected {
				t.Errorf("MCPConfigPath() = %q, want %q", got, tt.expected)
			}
		})
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
			name:              "cursor capabilities",
			wantOutputStyles:  false,
			wantSlashCommands: false,
			wantSkills:        true,
			wantSystemPrompt:  true,
			wantMCP:           true,
			wantAutoInstall:   false,
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
	// Cursor doesn't support slash commands
	if got := a.CommandsDir("/home/user"); got != "" {
		t.Errorf("CommandsDir() = %q, want empty string", got)
	}
}

func TestOutputStyleDir(t *testing.T) {
	a := NewAdapter()
	// Cursor doesn't support output styles
	if got := a.OutputStyleDir("/home/user"); got != "" {
		t.Errorf("OutputStyleDir() = %q, want empty string", got)
	}
}

func TestInstallCommandErrorType(t *testing.T) {
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
