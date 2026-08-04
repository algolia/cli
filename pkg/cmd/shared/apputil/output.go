package apputil

import (
	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
)

// ApplicationOutput is the machine-readable view of an application, shared by
// every command that emits one (`application current`, `application select`,
// `auth login`)
type ApplicationOutput struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
	Plan  string `json:"plan"`
}

// NewApplicationOutput builds the output view, reading the alias from the
// config so it reflects what was actually persisted.
func NewApplicationOutput(cfg config.IConfig, app *dashboard.Application) ApplicationOutput {
	out := ApplicationOutput{
		ID:   app.ID,
		Name: app.Name,
		Plan: app.PlanLabel,
	}
	if alias, ok := cfg.ApplicationAlias(app.ID); ok {
		out.Alias = alias
	}

	return out
}
