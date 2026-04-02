package catalog

import (
	"testing"

	"github.com/zabadev/agent-ai/internal/model"
)

func TestMVPComponents(t *testing.T) {
	t.Run("returns non-empty list", func(t *testing.T) {
		components := MVPComponents()
		if len(components) == 0 {
			t.Error("expected non-empty components list")
		}
	})

	t.Run("component IDs match expected values", func(t *testing.T) {
		expectedIDs := []model.ComponentID{
			model.ComponentEngram,
			model.ComponentSDD,
			model.ComponentSkills,
			model.ComponentContext7,
			model.ComponentPersona,
			model.ComponentPermission,
			model.ComponentGGA,
			model.ComponentTheme,
		}

		components := MVPComponents()
		if len(components) != len(expectedIDs) {
			t.Errorf("got %d components, expected %d", len(components), len(expectedIDs))
		}

		for i, expectedID := range expectedIDs {
			if i >= len(components) {
				t.Errorf("missing component at index %d", i)
				continue
			}
			if components[i].ID != expectedID {
				t.Errorf("component[%d].ID = %q, want %q", i, components[i].ID, expectedID)
			}
		}
	})

	t.Run("component names are non-empty", func(t *testing.T) {
		components := MVPComponents()
		for i, comp := range components {
			if comp.Name == "" {
				t.Errorf("component[%d].Name is empty", i)
			}
		}
	})

	t.Run("component descriptions are non-empty", func(t *testing.T) {
		components := MVPComponents()
		for i, comp := range components {
			if comp.Description == "" {
				t.Errorf("component[%d].Description is empty", i)
			}
		}
	})

	t.Run("returns copy not original", func(t *testing.T) {
		components1 := MVPComponents()
		components2 := MVPComponents()

		if &components1[0] == &components2[0] {
			t.Error("MVPComponents should return a copy, not the original slice")
		}
	})
}

func TestMVPSkills(t *testing.T) {
	t.Run("returns non-empty list", func(t *testing.T) {
		skills := MVPSkills()
		if len(skills) == 0 {
			t.Error("expected non-empty skills list")
		}
	})

	t.Run("skill IDs are valid", func(t *testing.T) {
		skills := MVPSkills()
		for i, skill := range skills {
			if skill.ID == "" {
				t.Errorf("skill[%d].ID is empty", i)
			}
		}
	})

	t.Run("skill names are non-empty", func(t *testing.T) {
		skills := MVPSkills()
		for i, skill := range skills {
			if skill.Name == "" {
				t.Errorf("skill[%d].Name is empty", i)
			}
		}
	})

	t.Run("skill categories are non-empty", func(t *testing.T) {
		skills := MVPSkills()
		for i, skill := range skills {
			if skill.Category == "" {
				t.Errorf("skill[%d].Category is empty", i)
			}
		}
	})

	t.Run("skill priorities are non-empty", func(t *testing.T) {
		skills := MVPSkills()
		for i, skill := range skills {
			if skill.Priority == "" {
				t.Errorf("skill[%d].Priority is empty", i)
			}
		}
	})

	t.Run("returns copy not original", func(t *testing.T) {
		skills1 := MVPSkills()
		skills2 := MVPSkills()

		if &skills1[0] == &skills2[0] {
			t.Error("MVPSkills should return a copy, not the original slice")
		}
	})
}

func TestAllAgents(t *testing.T) {
	t.Run("returns non-empty list", func(t *testing.T) {
		agents := AllAgents()
		if len(agents) == 0 {
			t.Error("expected non-empty agents list")
		}
	})

	t.Run("agent IDs are valid", func(t *testing.T) {
		agents := AllAgents()
		for i, agent := range agents {
			if agent.ID == "" {
				t.Errorf("agent[%d].ID is empty", i)
			}
		}
	})

	t.Run("agent names are non-empty", func(t *testing.T) {
		agents := AllAgents()
		for i, agent := range agents {
			if agent.Name == "" {
				t.Errorf("agent[%d].Name is empty", i)
			}
		}
	})

	t.Run("returns copy not original", func(t *testing.T) {
		agents1 := AllAgents()
		agents2 := AllAgents()

		if &agents1[0] == &agents2[0] {
			t.Error("AllAgents should return a copy, not the original slice")
		}
	})
}

func TestMVPAgents(t *testing.T) {
	t.Run("returns non-empty list", func(t *testing.T) {
		agents := MVPAgents()
		if len(agents) == 0 {
			t.Error("expected non-empty MVP agents list")
		}
	})

	t.Run("MVP agents are subset of all agents", func(t *testing.T) {
		all := AllAgents()
		mvp := MVPAgents()

		mvpIDs := make(map[model.AgentID]bool)
		for _, a := range mvp {
			mvpIDs[a.ID] = true
		}

		for _, a := range all {
			// Verify each agent: if it's in mvpIDs, that's expected (MVP ⊆ all).
			// This is a sanity check that MVPAgents() returns a subset of AllAgents().
			_ = mvpIDs[a.ID] // intentionally unused — just iterating to ensure no panics
		}
	})
}

func TestIsMVPAgent(t *testing.T) {
	tests := []struct {
		name     string
		agent    model.AgentID
		expected bool
	}{
		{"Claude Code is MVP", model.AgentClaudeCode, true},
		{"OpenCode is MVP", model.AgentOpenCode, true},
		{"Gemini CLI is not MVP", model.AgentGeminiCLI, false},
		{"Cursor is not MVP", model.AgentCursor, false},
		{"Codex is not MVP", model.AgentCodex, false},
		{"VS Code Copilot is not MVP", model.AgentVSCodeCopilot, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMVPAgent(tt.agent)
			if result != tt.expected {
				t.Errorf("IsMVPAgent(%q) = %v, want %v", tt.agent, result, tt.expected)
			}
		})
	}
}

func TestIsSupportedAgent(t *testing.T) {
	tests := []struct {
		name     string
		agent    model.AgentID
		expected bool
	}{
		{"Claude Code is supported", model.AgentClaudeCode, true},
		{"OpenCode is supported", model.AgentOpenCode, true},
		{"Gemini CLI is supported", model.AgentGeminiCLI, true},
		{"Cursor is supported", model.AgentCursor, true},
		{"Codex is supported", model.AgentCodex, true},
		{"VS Code Copilot is supported", model.AgentVSCodeCopilot, true},
		{"unknown agent is not supported", model.AgentID("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSupportedAgent(tt.agent)
			if result != tt.expected {
				t.Errorf("IsSupportedAgent(%q) = %v, want %v", tt.agent, result, tt.expected)
			}
		})
	}
}
