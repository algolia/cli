// This file is generated; DO NOT EDIT.

package cmdutil

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var AgentCompletionRequest = []string{
	"algolia",
	"configuration",
	"id",
	"messages",
	"toolApprovals",
}

var AgentConfigCreate = []string{
	"config",
	"description",
	"instructions",
	"model",
	"name",
	"providerId",
	"systemPrompt",
	"templateType",
	"tools",
}

func AddAgentCompletionRequestFlags(cmd *cobra.Command) {
	algolia := NewJSONVar([]string{}...)
	cmd.Flags().Var(algolia, "algolia", heredoc.Doc(`.`))
	configuration := NewJSONVar([]string{}...)
	cmd.Flags().Var(configuration, "configuration", heredoc.Doc(`Dynamic configuration for testing agents.`))
	cmd.Flags().String("id", "", heredoc.Doc(`.`))
	messages := NewJSONVar([]string{"array", "array"}...)
	cmd.Flags().Var(messages, "messages", heredoc.Doc(`.`))
	toolApprovals := NewJSONVar([]string{}...)
	cmd.Flags().Var(toolApprovals, "toolApprovals", heredoc.Doc(`.`))
}

func AddAgentConfigCreateFlags(cmd *cobra.Command) {
	config := NewJSONVar([]string{}...)
	cmd.Flags().Var(config, "config", heredoc.Doc(`.`))
	cmd.Flags().String("description", "", heredoc.Doc(`.`))
	cmd.Flags().String("instructions", "", heredoc.Doc(`The agent prompt: defines the agent's role, tone, and goals. Guides how it answers using the provided context. Corresponds to the 'Agent prompt' field in the dashboard.`))
	cmd.Flags().String("model", "", heredoc.Doc(`.`))
	cmd.Flags().String("name", "", heredoc.Doc(`.`))
	cmd.Flags().String("providerId", "", heredoc.Doc(`.`))
	cmd.Flags().String("systemPrompt", "", heredoc.Doc(`.`))
	cmd.Flags().String("templateType", "", heredoc.Doc(`.`))
	tools := NewJSONVar([]string{}...)
	cmd.Flags().Var(tools, "tools", heredoc.Doc(`.`))
}
