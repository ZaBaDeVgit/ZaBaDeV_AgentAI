package opencode

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
			lookPathPath:    "/opt/homebrew/bin/opencode",
			stat:            statResult{isDir: true},
			wantInstalled:   true,
			wantBinaryPath:  "/opt/homebrew/bin/opencode",
			wantConfigPath:  filepath.Join("/tmp/home", ".config", "opencode"),
			wantConfigFound: true,
		},
		{
			name:            "binary missing and config missing",
			lookPathErr:     errors.New("missing"),
			stat:            statResult{err: os.ErrNotExist},
			wantInstalled:   false,
			wantBinaryPath:  "",
			wantConfigPath:  filepath.Join("/tmp/home", ".config", "opencode"),
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
		wantErr bool
	}{
		{
			name:    "darwin resolves official anomalyco brew tap",
			profile: system.PlatformProfile{OS: "darwin", PackageManager: "brew"},
			want:    [][]string{{"brew", "install", "anomalyco/tap/opencode"}},
		},
		{
			name:    "ubuntu resolves npm install",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroUbuntu, PackageManager: "apt"},
			want:    [][]string{{"sudo", "npm", "install", "-g", "opencode-ai"}},
		},
		{
			name:    "arch resolves npm install",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroArch, PackageManager: "pacman"},
			want:    [][]string{{"sudo", "npm", "install", "-g", "opencode-ai"}},
		},
		{
			name:    "fedora resolves npm install",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroFedora, PackageManager: "dnf"},
			want:    [][]string{{"sudo", "npm", "install", "-g", "opencode-ai"}},
		},
		{
			name:    "fedora with writable npm skips sudo",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroFedora, PackageManager: "dnf", NpmWritable: true},
			want:    [][]string{{"npm", "install", "-g", "opencode-ai"}},
		},
		{
			name:    "unsupported package manager returns error",
			profile: system.PlatformProfile{OS: "linux", LinuxDistro: system.LinuxDistroUbuntu, PackageManager: "zypper"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := a.InstallCommand(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("InstallCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(command, tt.want) {
				t.Fatalf("InstallCommand() = %v, want %v", command, tt.want)
			}
		})
	}
}

func TestAdapterIdentity(t *testing.T) {
	tests := []struct {
		name      string
		wantAgent model.AgentID
		wantTier  model.SupportTier
	}{
		{
			name:      "opencode agent identity",
			wantAgent: model.AgentOpenCode,
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
			expected: "/home/user/.config/opencode/AGENTS.md",
		},
		{
			name:     "path with special characters",
			homeDir:  "/home/user-name_123",
			expected: "/home/user-name_123/.config/opencode/AGENTS.md",
		},
		{
			name:     "empty home dir returns relative path",
			homeDir:  "",
			expected: ".config/opencode/AGENTS.md",
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
			expected: "/home/user/.config/opencode",
		},
		{
			name:     "empty home returns relative path",
			homeDir:  "",
			expected: ".config/opencode",
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
			expected: "/home/user/.config/opencode",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.config/opencode",
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
			expected: "/home/user/.config/opencode/skills",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.config/opencode/skills",
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
			expected: "/home/user/.config/opencode/opencode.json",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.config/opencode/opencode.json",
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

func TestStrategies(t *testing.T) {
	tests := []struct {
		name                     string
		wantSystemPromptStrategy model.SystemPromptStrategy
		wantMCPStrategy          model.MCPStrategy
	}{
		{
			name:                     "opencode strategies",
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

func TestMCPConfigPath(t *testing.T) {
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
			expected:   "/home/user/.config/opencode/opencode.json",
		},
		{
			name:       "filesystem server",
			homeDir:    "/home/user",
			serverName: "filesystem",
			expected:   "/home/user/.config/opencode/opencode.json",
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
			name:              "opencode capabilities",
			wantOutputStyles:  false,
			wantSlashCommands: true,
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
	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home",
			homeDir:  "/home/user",
			expected: "/home/user/.config/opencode/commands",
		},
		{
			name:     "tmp directory",
			homeDir:  "/tmp/home",
			expected: "/tmp/home/.config/opencode/commands",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.CommandsDir(tt.homeDir)
			if got != tt.expected {
				t.Errorf("CommandsDir() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestOutputStyleDir(t *testing.T) {
	a := NewAdapter()
	// OpenCode doesn't support output styles
	if got := a.OutputStyleDir("/home/user"); got != "" {
		t.Errorf("OutputStyleDir() = %q, want empty string", got)
	}
}
