package gemini

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/zabadev/agent-ai/internal/model"
	"github.com/zabadev/agent-ai/internal/system"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name            string
		lookPathPath    string
		lookPathErr     error
		stat            statResult
		wantInstalled   bool
		wantBinaryPath  string
		wantConfigPath  string
		wantConfigFound bool
		wantErr         bool
	}{
		{
			name:            "binary and config directory found",
			lookPathPath:    "/usr/local/bin/gemini",
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantBinaryPath:  "/usr/local/bin/gemini",
			wantConfigPath:  filepath.Join("/tmp/home", ".gemini"),
			wantConfigFound: true,
		},
		{
			name:            "binary missing and config missing",
			lookPathErr:     errors.New("missing"),
			stat:            statResult{err: os.ErrNotExist},
			wantInstalled:   false,
			wantBinaryPath:  "",
			wantConfigPath:  filepath.Join("/tmp/home", ".gemini"),
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
				lookPath: func(string) (string, error) {
					return tt.lookPathPath, tt.lookPathErr
				},
				statPath: func(string) statResult {
					return tt.stat
				},
			}

			installed, binaryPath, configPath, configFound, err := a.Detect(context.Background(), "/tmp/home")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Detect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if installed != tt.wantInstalled {
				t.Fatalf("Detect() installed = %v, want %v", installed, tt.wantInstalled)
			}

			if binaryPath != tt.wantBinaryPath {
				t.Fatalf("Detect() binaryPath = %q, want %q", binaryPath, tt.wantBinaryPath)
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

func TestInstallCommand(t *testing.T) {
	a := NewAdapter()

	tests := []struct {
		name    string
		profile system.PlatformProfile
		want    [][]string
	}{
		{
			name:    "darwin uses npm without sudo",
			profile: system.PlatformProfile{OS: "darwin", PackageManager: "brew"},
			want:    [][]string{{"npm", "install", "-g", "@google/gemini-cli"}},
		},
		{
			name:    "linux system npm uses sudo",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroUbuntu, PackageManager: "apt"},
			want:    [][]string{{"sudo", "npm", "install", "-g", "@google/gemini-cli"}},
		},
		{
			name:    "linux nvm skips sudo",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroUbuntu, PackageManager: "apt", NpmWritable: true},
			want:    [][]string{{"npm", "install", "-g", "@google/gemini-cli"}},
		},
		{
			name:    "windows uses npm without sudo",
			profile: system.PlatformProfile{OS: "windows", PackageManager: "winget", NpmWritable: true},
			want:    [][]string{{"npm", "install", "-g", "@google/gemini-cli"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := a.InstallCommand(tt.profile)
			if err != nil {
				t.Fatalf("InstallCommand() returned error: %v", err)
			}

			if !reflect.DeepEqual(command, tt.want) {
				t.Fatalf("InstallCommand() = %v, want %v", command, tt.want)
			}
		})
	}
}

func TestConfigPathsCrossPlatform(t *testing.T) {
	a := NewAdapter()
	home := "/tmp/home"

	if got := a.GlobalConfigDir(home); got != filepath.Join(home, ".gemini") {
		t.Fatalf("GlobalConfigDir() = %q, want %q", got, filepath.Join(home, ".gemini"))
	}

	if got := a.SkillsDir(home); got != filepath.Join(home, ".gemini", "skills") {
		t.Fatalf("SkillsDir() = %q, want %q", got, filepath.Join(home, ".gemini", "skills"))
	}

	if got := a.MCPConfigPath(home, "ctx7"); got != filepath.Join(home, ".gemini", "settings.json") {
		t.Fatalf("MCPConfigPath() = %q, want %q", got, filepath.Join(home, ".gemini", "settings.json"))
	}

	if got := a.SystemPromptFile(home); got != filepath.Join(home, ".gemini", "GEMINI.md") {
		t.Fatalf("SystemPromptFile() = %q, want %q", got, filepath.Join(home, ".gemini", "GEMINI.md"))
	}
}

func TestAdapterIdentity(t *testing.T) {
	tests := []struct {
		name      string
		wantAgent model.AgentID
		wantTier  model.SupportTier
	}{
		{
			name:      "gemini agent identity",
			wantAgent: model.AgentGeminiCLI,
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
			expected: "/home/user/.gemini/GEMINI.md",
		},
		{
			name:     "path with special characters",
			homeDir:  "/home/user-name_123",
			expected: "/home/user-name_123/.gemini/GEMINI.md",
		},
		{
			name:     "empty home dir returns relative path",
			homeDir:  "",
			expected: ".gemini/GEMINI.md",
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
			expected: "/home/user/.gemini",
		},
		{
			name:     "empty home returns relative path",
			homeDir:  "",
			expected: ".gemini",
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
			expected: "/home/user/.gemini",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.gemini",
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
			expected: "/home/user/.gemini/skills",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.gemini/skills",
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

func TestSettingsPath(t *testing.T) {
	a := NewAdapter()
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.gemini/settings.json",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.gemini/settings.json",
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

func TestStrategiesTableDriven(t *testing.T) {
	tests := []struct {
		name                     string
		wantSystemPromptStrategy model.SystemPromptStrategy
		wantMCPStrategy          model.MCPStrategy
	}{
		{
			name:                     "gemini strategies",
			wantSystemPromptStrategy: model.StrategyFileReplace,
			wantMCPStrategy:          model.StrategyMergeIntoSettings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter()
			if got := a.SystemPromptStrategy(); got != tt.wantSystemPromptStrategy {
				t.Errorf("SystemPromptStrategy() = %v, want %v", got, tt.wantSystemPromptStrategy)
			}
			if got := a.MCPStrategy(); got != tt.wantMCPStrategy {
				t.Errorf("MCPStrategy() = %v, want %v", got, tt.wantMCPStrategy)
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
			expected:   "/home/user/.gemini/settings.json",
		},
		{
			name:       "filesystem server",
			homeDir:    "/home/user",
			serverName: "filesystem",
			expected:   "/home/user/.gemini/settings.json",
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
			name:              "gemini capabilities",
			wantOutputStyles:  false,
			wantSlashCommands: false,
			wantSkills:        true,
			wantSystemPrompt:  true,
			wantMCP:           true,
			wantAutoInstall:   true,
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
	// Gemini doesn't support slash commands
	if got := a.CommandsDir("/home/user"); got != "" {
		t.Errorf("CommandsDir() = %q, want empty string", got)
	}
}

func TestOutputStyleDir(t *testing.T) {
	a := NewAdapter()
	// Gemini doesn't support output styles
	if got := a.OutputStyleDir("/home/user"); got != "" {
		t.Errorf("OutputStyleDir() = %q, want empty string", got)
	}
}
