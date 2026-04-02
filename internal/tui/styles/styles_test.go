package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderLogo(t *testing.T) {
	t.Run("returns non-empty string", func(t *testing.T) {
		logo := RenderLogo()
		if logo == "" {
			t.Error("RenderLogo() returned empty string")
		}
	})

	t.Run("returns consistent output", func(t *testing.T) {
		logo1 := RenderLogo()
		logo2 := RenderLogo()
		if logo1 != logo2 {
			t.Error("RenderLogo() should return consistent output")
		}
	})

	t.Run("contains expected ASCII art characters", func(t *testing.T) {
		logo := RenderLogo()
		// Logo should contain box-drawing characters from the ASCII art
		hasBoxChar := false
		for _, ch := range logo {
			if ch == '█' || ch == '╗' || ch == '╝' || ch == '╔' || ch == '╚' {
				hasBoxChar = true
				break
			}
		}
		if !hasBoxChar {
			t.Error("RenderLogo() should contain box-drawing characters")
		}
	})

	t.Run("has correct number of lines", func(t *testing.T) {
		logo := RenderLogo()
		lineCount := 0
		for _, ch := range logo {
			if ch == '\n' {
				lineCount++
			}
		}
		// 6 lines in logoLines means 5 newlines between them
		if lineCount != 5 {
			t.Errorf("RenderLogo() should have 5 newlines, got %d", lineCount)
		}
	})
}

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name  string
		value lipgloss.Color
	}{
		{"ColorBase", ColorBase},
		{"ColorSurface", ColorSurface},
		{"ColorOverlay", ColorOverlay},
		{"ColorText", ColorText},
		{"ColorSubtext", ColorSubtext},
		{"ColorLavender", ColorLavender},
		{"ColorGreen", ColorGreen},
		{"ColorPeach", ColorPeach},
		{"ColorRed", ColorRed},
		{"ColorBlue", ColorBlue},
		{"ColorMauve", ColorMauve},
		{"ColorYellow", ColorYellow},
		{"ColorTeal", ColorTeal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("color constant %s is empty", tt.name)
			}
		})
	}
}

func TestTagline(t *testing.T) {
	tests := []struct {
		version  string
		expected string
	}{
		{"v1.0.0", "ZaBaDeV v1.0.0 — One command. Any agent. Any OS."},
		{"v2.0.0-beta", "ZaBaDeV v2.0.0-beta — One command. Any agent. Any OS."},
		{"dev", "ZaBaDeV dev — One command. Any agent. Any OS."},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := Tagline(tt.version)
			if result != tt.expected {
				t.Errorf("Tagline(%q) = %q, want %q", tt.version, result, tt.expected)
			}
		})
	}
}

func TestCursor(t *testing.T) {
	if Cursor == "" {
		t.Error("Cursor constant should not be empty")
	}
	if Cursor != "▸ " {
		t.Errorf("Cursor = %q, want %q", Cursor, "▸ ")
	}
}

func TestStyleConstants(t *testing.T) {
	// Verify that all style constants are non-nil by checking they can render
	// Since lipgloss.Style is a value type, we verify styles can be applied

	testInput := "test content"

	styles := map[string]lipgloss.Style{
		"TitleStyle":      TitleStyle,
		"HeadingStyle":    HeadingStyle,
		"HelpStyle":       HelpStyle,
		"SubtextStyle":    SubtextStyle,
		"SelectedStyle":   SelectedStyle,
		"UnselectedStyle": UnselectedStyle,
		"SuccessStyle":    SuccessStyle,
		"ErrorStyle":      ErrorStyle,
		"WarningStyle":    WarningStyle,
		"FrameStyle":      FrameStyle,
		"PanelStyle":      PanelStyle,
		"ProgressFilled":  ProgressFilled,
		"ProgressEmpty":   ProgressEmpty,
		"PercentStyle":    PercentStyle,
	}

	for name, style := range styles {
		t.Run(name, func(t *testing.T) {
			result := style.Render(testInput)
			if result == "" {
				t.Errorf("%s.Render() returned empty string", name)
			}
		})
	}
}
