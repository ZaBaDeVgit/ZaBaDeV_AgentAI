package model

import "testing"

func TestAgentIDConstants(t *testing.T) {
	tests := []struct {
		name  string
		value AgentID
	}{
		{"AgentClaudeCode", AgentClaudeCode},
		{"AgentOpenCode", AgentOpenCode},
		{"AgentGeminiCLI", AgentGeminiCLI},
		{"AgentCursor", AgentCursor},
		{"AgentVSCodeCopilot", AgentVSCodeCopilot},
		{"AgentCodex", AgentCodex},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("AgentID %s is empty", tt.name)
			}
		})
	}

	t.Run("all values are unique", func(t *testing.T) {
		seen := map[AgentID]bool{}
		for _, tt := range tests {
			if seen[tt.value] {
				t.Errorf("duplicate AgentID value: %s", tt.value)
			}
			seen[tt.value] = true
		}
	})
}

func TestComponentIDConstants(t *testing.T) {
	tests := []struct {
		name  string
		value ComponentID
	}{
		{"ComponentEngram", ComponentEngram},
		{"ComponentSDD", ComponentSDD},
		{"ComponentSkills", ComponentSkills},
		{"ComponentContext7", ComponentContext7},
		{"ComponentPersona", ComponentPersona},
		{"ComponentPermission", ComponentPermission},
		{"ComponentGGA", ComponentGGA},
		{"ComponentTheme", ComponentTheme},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("ComponentID %s is empty", tt.name)
			}
		})
	}

	t.Run("all values are unique", func(t *testing.T) {
		seen := map[ComponentID]bool{}
		for _, tt := range tests {
			if seen[tt.value] {
				t.Errorf("duplicate ComponentID value: %s", tt.value)
			}
			seen[tt.value] = true
		}
	})
}

func TestSkillIDConstants(t *testing.T) {
	// Test a subset of skill IDs to ensure they are properly defined
	tests := []struct {
		name  string
		value SkillID
	}{
		{"SkillSDDInit", SkillSDDInit},
		{"SkillSDDApply", SkillSDDApply},
		{"SkillSDDVerify", SkillSDDVerify},
		{"SkillSDDExplore", SkillSDDExplore},
		{"SkillSDDPropose", SkillSDDPropose},
		{"SkillSDDSpec", SkillSDDSpec},
		{"SkillSDDDesign", SkillSDDDesign},
		{"SkillSDDTasks", SkillSDDTasks},
		{"SkillSDDArchive", SkillSDDArchive},
		{"SkillGoTesting", SkillGoTesting},
		{"SkillCreator", SkillCreator},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("SkillID %s is empty", tt.name)
			}
		})
	}
}

func TestPersonaIDConstants(t *testing.T) {
	tests := []struct {
		name  string
		value PersonaID
	}{
		{"PersonaGentleman", PersonaGentleman},
		{"PersonaSeniorZaBaDeV", PersonaSeniorZaBaDeV},
		{"PersonaNeutral", PersonaNeutral},
		{"PersonaCustom", PersonaCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("PersonaID %s is empty", tt.name)
			}
		})
	}
}

func TestPresetIDConstants(t *testing.T) {
	tests := []struct {
		name  string
		value PresetID
	}{
		{"PresetFullGentleman", PresetFullGentleman},
		{"PresetEcosystemOnly", PresetEcosystemOnly},
		{"PresetMinimal", PresetMinimal},
		{"PresetCustom", PresetCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("PresetID %s is empty", tt.name)
			}
		})
	}
}

func TestSDDModeIDConstants(t *testing.T) {
	tests := []struct {
		name  string
		value SDDModeID
	}{
		{"SDDModeSingle", SDDModeSingle},
		{"SDDModeMulti", SDDModeMulti},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("SDDModeID %s is empty", tt.name)
			}
		})
	}
}

func TestSupportTier(t *testing.T) {
	if TierFull == "" {
		t.Error("SupportTier TierFull is empty")
	}
}

func TestSystemPromptStrategy(t *testing.T) {
	tests := []struct {
		name  string
		value SystemPromptStrategy
	}{
		{"StrategyMarkdownSections", StrategyMarkdownSections},
		{"StrategyFileReplace", StrategyFileReplace},
		{"StrategyAppendToFile", StrategyAppendToFile},
		{"StrategyInstructionsFile", StrategyInstructionsFile},
	}

	// Verify all strategy values are defined and distinct
	t.Run("all values are distinct", func(t *testing.T) {
		seen := map[SystemPromptStrategy]bool{}
		for _, tt := range tests {
			if seen[tt.value] {
				t.Errorf("duplicate SystemPromptStrategy value: %s", tt.name)
			}
			seen[tt.value] = true
		}
		if len(seen) != len(tests) {
			t.Errorf("expected %d unique strategies, got %d", len(tests), len(seen))
		}
	})

	t.Run("values are sequential iota", func(t *testing.T) {
		if StrategyFileReplace != StrategyMarkdownSections+1 {
			t.Errorf("StrategyFileReplace should be StrategyMarkdownSections+1")
		}
		if StrategyAppendToFile != StrategyFileReplace+1 {
			t.Errorf("StrategyAppendToFile should be StrategyFileReplace+1")
		}
		if StrategyInstructionsFile != StrategyAppendToFile+1 {
			t.Errorf("StrategyInstructionsFile should be StrategyAppendToFile+1")
		}
	})
}

func TestMCPStrategy(t *testing.T) {
	tests := []struct {
		name  string
		value MCPStrategy
	}{
		{"StrategySeparateMCPFiles", StrategySeparateMCPFiles},
		{"StrategyMergeIntoSettings", StrategyMergeIntoSettings},
		{"StrategyMCPConfigFile", StrategyMCPConfigFile},
		{"StrategyTOMLFile", StrategyTOMLFile},
	}

	// Verify all strategy values are defined and distinct
	t.Run("all values are distinct", func(t *testing.T) {
		seen := map[MCPStrategy]bool{}
		for _, tt := range tests {
			if seen[tt.value] {
				t.Errorf("duplicate MCPStrategy value: %s", tt.name)
			}
			seen[tt.value] = true
		}
		if len(seen) != len(tests) {
			t.Errorf("expected %d unique strategies, got %d", len(tests), len(seen))
		}
	})

	t.Run("values are sequential iota", func(t *testing.T) {
		if StrategyMergeIntoSettings != StrategySeparateMCPFiles+1 {
			t.Errorf("StrategyMergeIntoSettings should be StrategySeparateMCPFiles+1")
		}
		if StrategyMCPConfigFile != StrategyMergeIntoSettings+1 {
			t.Errorf("StrategyMCPConfigFile should be StrategyMergeIntoSettings+1")
		}
		if StrategyTOMLFile != StrategyMCPConfigFile+1 {
			t.Errorf("StrategyTOMLFile should be StrategyMCPConfigFile+1")
		}
	})
}

func TestPlanStatus(t *testing.T) {
	tests := []struct {
		name  string
		value PlanStatus
	}{
		{"PlanStatusPending", PlanStatusPending},
		{"PlanStatusRunning", PlanStatusRunning},
		{"PlanStatusSucceeded", PlanStatusSucceeded},
		{"PlanStatusFailed", PlanStatusFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("PlanStatus %s is empty", tt.name)
			}
		})
	}
}

func TestRunResult(t *testing.T) {
	tests := []struct {
		name  string
		value RunResult
	}{
		{"RunResultSkipped", RunResultSkipped},
		{"RunResultSuccess", RunResultSuccess},
		{"RunResultFailed", RunResultFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("RunResult %s is empty", tt.name)
			}
		})
	}
}

func TestModelAssignmentFullID(t *testing.T) {
	tests := []struct {
		name         string
		providerID   string
		modelID      string
		expectedFull string
	}{
		{"anthropic claude", "anthropic", "claude-sonnet-4-20250514", "anthropic/claude-sonnet-4-20250514"},
		{"openai gpt", "openai", "gpt-4o", "openai/gpt-4o"},
		{"google gemini", "google", "gemini-2.0-flash", "google/gemini-2.0-flash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ModelAssignment{
				ProviderID: tt.providerID,
				ModelID:    tt.modelID,
			}
			result := m.FullID()
			if result != tt.expectedFull {
				t.Errorf("FullID() = %q, want %q", result, tt.expectedFull)
			}
		})
	}
}

func TestSelectionHasAgent(t *testing.T) {
	s := Selection{
		Agents: []AgentID{AgentClaudeCode, AgentOpenCode},
	}

	tests := []struct {
		name     string
		agent    AgentID
		expected bool
	}{
		{"Claude Code present", AgentClaudeCode, true},
		{"OpenCode present", AgentOpenCode, true},
		{"Gemini CLI absent", AgentGeminiCLI, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.HasAgent(tt.agent)
			if result != tt.expected {
				t.Errorf("HasAgent(%q) = %v, want %v", tt.agent, result, tt.expected)
			}
		})
	}
}

func TestSelectionHasComponent(t *testing.T) {
	s := Selection{
		Components: []ComponentID{ComponentEngram, ComponentSDD},
	}

	tests := []struct {
		name      string
		component ComponentID
		expected  bool
	}{
		{"Engram present", ComponentEngram, true},
		{"SDD present", ComponentSDD, true},
		{"Skills absent", ComponentSkills, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.HasComponent(tt.component)
			if result != tt.expected {
				t.Errorf("HasComponent(%q) = %v, want %v", tt.component, result, tt.expected)
			}
		})
	}
}
