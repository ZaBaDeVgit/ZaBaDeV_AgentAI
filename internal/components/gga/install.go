package gga

import (
	"github.com/zabadev/agent-ai/internal/installcmd"
	"github.com/zabadev/agent-ai/internal/model"
	"github.com/zabadev/agent-ai/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentGGA)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
