package catalog

import "github.com/zabadev/agent-ai/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// SDD skills (10 phases)
	{ID: model.SkillSDDInit, Name: "sdd-init", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDNew, Name: "sdd-new", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDApply, Name: "sdd-apply", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDVerify, Name: "sdd-verify", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDExplore, Name: "sdd-explore", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDPropose, Name: "sdd-propose", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDSpec, Name: "sdd-spec", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDDesign, Name: "sdd-design", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDTasks, Name: "sdd-tasks", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDArchive, Name: "sdd-archive", Category: "sdd", Priority: "p0"},
	// Framework/Language skills
	{ID: model.SkillReact19, Name: "react-19", Category: "frontend", Priority: "p0"},
	{ID: model.SkillNextJS15, Name: "nextjs-15", Category: "frontend", Priority: "p0"},
	{ID: model.SkillTailwind4, Name: "tailwind-4", Category: "styling", Priority: "p0"},
	{ID: model.SkillZustand5, Name: "zustand-5", Category: "state", Priority: "p0"},
	{ID: model.SkillZod4, Name: "zod-4", Category: "validation", Priority: "p0"},
	{ID: model.SkillAISDK5, Name: "ai-sdk-5", Category: "ai", Priority: "p0"},
	// Testing skills
	{ID: model.SkillPlaywright, Name: "playwright", Category: "testing", Priority: "p0"},
	{ID: model.SkillPytest, Name: "pytest", Category: "testing", Priority: "p0"},
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	// Backend skills
	{ID: model.SkillDjangoDRF, Name: "django-drf", Category: "backend", Priority: "p0"},
	// Workflow/Builder skills
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
