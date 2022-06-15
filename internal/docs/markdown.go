package docs

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func printOptions(w io.Writer, cmd *cobra.Command) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(w)
	if flags.HasAvailableFlags() {
		fmt.Fprint(w, "### Options\n\n")
		if err := printFlagsHTML(w, flags); err != nil {
			return err
		}
		fmt.Fprint(w, "\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(w)
	if hasNonHelpFlags(parentFlags) {
		fmt.Fprint(w, "### Options inherited from parent commands\n\n")
		if err := printFlagsHTML(w, parentFlags); err != nil {
			return err
		}
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

func hasNonHelpFlags(fs *pflag.FlagSet) (found bool) {
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden && f.Name != "help" {
			found = true
		}
	})
	return
}

type flagView struct {
	Name      string
	Varname   string
	Shorthand string
	Usage     string
}

var flagsTemplate = `
<dl class="flags">{{ range . }}
	<dt>{{ if .Shorthand }}<code>-{{.Shorthand}}</code>, {{ end -}}
		<code>--{{.Name}}{{ if .Varname }} &lt;{{.Varname}}&gt;{{ end }}</code></dt>
	<dd>{{.Usage}}</dd>
{{ end }}</dl>
`

var tpl = template.Must(template.New("flags").Parse(flagsTemplate))

func printFlagsHTML(w io.Writer, fs *pflag.FlagSet) error {
	var flags []flagView
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		varname, usage := pflag.UnquoteUsage(f)
		flags = append(flags, flagView{
			Name:      f.Name,
			Varname:   varname,
			Shorthand: f.Shorthand,
			Usage:     usage,
		})
	})
	return tpl.Execute(w, flags)
}

type CmdView struct {
	Title       string
	Parent      string
	HasChildren bool
}

var headerTemplate = `---
layout: default
title: {{.Title}}{{if .HasChildren}}
has_children: true
{{end}}
{{if .Parent}}parent: {{.Parent}}{{end}}
---
`

var headerTpl = template.Must(template.New("header").Parse(headerTemplate))

func printHeaderMarkdown(w io.Writer, cmd *cobra.Command) error {
	// Do no put everything under the "algolia" command
	if cmd.Name() == "algolia" {
		return headerTpl.Execute(w, &CmdView{
			Title:       cmd.CommandPath(),
			Parent:      "",
			HasChildren: false,
		})
	}

	var parent string
	if cmd.HasParent() {
		parent = cmd.Parent().CommandPath()
		if parent == "algolia" {
			parent = ""
		}
	}

	return headerTpl.Execute(w, &CmdView{
		Title:       cmd.CommandPath(),
		Parent:      parent,
		HasChildren: cmd.HasSubCommands(),
	})
}

func GenMarkdown(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	// Markdown header (with parent page if applicable)
	if err := printHeaderMarkdown(w, cmd); err != nil {
		return err
	}

	fmt.Fprintf(w, "## %s\n\n", cmd.CommandPath())

	hasLong := cmd.Long != ""
	if !hasLong {
		fmt.Fprintf(w, "%s\n\n", cmd.Short)
	}
	if cmd.Runnable() {
		fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine())
	}
	if hasLong {
		fmt.Fprintf(w, "%s\n\n", cmd.Long)
	}

	if cmd.Commands() != nil && len(cmd.Commands()) > 0 {
		fmt.Fprint(w, "### Commands\n\n")
		for _, subcmd := range cmd.Commands() {
			if !subcmd.IsAvailableCommand() {
				continue
			}
			fmt.Fprintf(w, "* [%s](%s)\n", subcmd.CommandPath(), linkHandler(cmdDocsPath(subcmd)))
		}
	}

	if err := printOptions(w, cmd); err != nil {
		return err
	}

	if len(cmd.Example) > 0 {
		fmt.Fprint(w, "### Examples\n\n{% highlight bash %}{% raw %}\n")
		fmt.Fprint(w, cmd.Example)
		fmt.Fprint(w, "{% endraw %}{% endhighlight %}\n\n")
	}

	if cmd.HasParent() {
		p := cmd.Parent()
		fmt.Fprint(w, "### See also\n\n")
		fmt.Fprintf(w, "* [%s](%s)\n", p.CommandPath(), linkHandler(cmdDocsPath(p)))
	}

	return nil
}

func GenMarkdownTree(cmd *cobra.Command, dir string, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}

		if err := GenMarkdownTree(c, dir, linkHandler); err != nil {
			return err
		}
	}

	filename := filepath.Join(dir, cmdDocsPath(cmd))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := GenMarkdown(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func cmdDocsPath(c *cobra.Command) string {
	return strings.ReplaceAll(c.CommandPath(), " ", "_") + ".md"
}
