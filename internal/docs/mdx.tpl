---
navigation: "cli"
title: |-
{{ .Name }}
description: |-
{{ .Description }}
slug: tools/cli/commands/{{ .Slug }}
---
{{ if .SubCommands }}
{{ range $subCommand := .SubCommands }}
## {{ $subCommand.Name }}

{{ $subCommand.Description }}

### Usage

`{{ $subCommand.Usage }}`

{{ if $subCommand.Examples }}### Examples{{ end }}
{{ $subCommand.Examples }}
{{ range $flagKey, $flagSlice := $subCommand.Flags }}
{{ if $flagSlice }}### Flags {{ end }}
{{ range $flag := $flagSlice }}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}
{{ end }}
{{ end }}
{{ else }}
## {{ .Name }}
### Usage

`{{ .Usage }}`

{{ if .Examples }}
### Examples
{{ .Examples }}
{{ end }}
{{ range $flagKey, $flagSlice := .Flags }}
### {{ $flagKey }}
{{ range $flag := $flagSlice }}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}
{{ end }}
{{end}}