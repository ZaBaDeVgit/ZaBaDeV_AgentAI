package screens

import (
	"strings"

	"github.com/zabadev/agent-ai/internal/model"
	"github.com/zabadev/agent-ai/internal/tui/styles"
)

func PersonaOptions() []model.PersonaID {
	return []model.PersonaID{
		model.PersonaSeniorZaBaDeV,
		model.PersonaNeutral,
		model.PersonaCustom,
	}
}

var personaLabels = map[model.PersonaID]string{
	model.PersonaSeniorZaBaDeV: "senior-zabadev — principal engineer",
	model.PersonaNeutral:       "neutral — direct and minimal",
	model.PersonaCustom:        "custom — keep my own prompt",
}

var personaDescriptions = map[model.PersonaID]string{
	model.PersonaSeniorZaBaDeV: "Single public agent. Installs Senior ZaBaDeV as the primary, memory-enabled, SDD-capable agent.",
	model.PersonaNeutral:       "No personality, just clean technical guidance.",
	model.PersonaCustom:        "Do not overwrite your existing persona instructions.",
}

func RenderPersona(selected model.PersonaID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Choose your Persona"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Pick the default behavior Gentle-AI will install for your selected agents."))
	b.WriteString("\n\n")

	for idx, persona := range PersonaOptions() {
		isSelected := persona == selected
		focused := idx == cursor
		b.WriteString(renderRadio(personaLabels[persona], isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+personaDescriptions[persona]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Back"}, cursor-len(PersonaOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return b.String()
}
