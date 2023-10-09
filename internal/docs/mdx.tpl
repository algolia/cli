---
navigation: "cli"
title: |-
{{ .Name }}
description: |-
{{ .Description }}
slug: tools/cli/commands/{{ .Slug }}
---
{{ if .SubCommands }}{{ range $subCommand := .SubCommands }}{{ if $subCommand.SubCommands }}{{ range $susCommand := $subCommand.SubCommands}}
## {{ $susCommand.Name }}

{{ $susCommand.Description }}

### Usage

`{{ $susCommand.Usage }}`

{{ $examples := getExamples $susCommand }}{{ if $examples }}
### Examples
{{ range $example := $examples }}{{ $example.Desc }}

```sh 
{{ if $example.WebCLICommand }}command="{{$example.WebCLICommand}}"{{ end }}
{{ $example.Code }}
```
{{ end }}{{ end }}{{ range $flagKey, $flagSlice := $susCommand.Flags }}{{ if $flagSlice }}
### Flags 
{{ range $flag := $flagSlice }}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}{{ end }}{{ end }}{{ end }}
{{ else }}
## {{ $subCommand.Name }}

{{ $subCommand.Description }}

### Usage

`{{ $subCommand.Usage }}`

{{ $examples := getExamples $subCommand }}
{{ if $examples }}
### Examples
{{ range $example := $examples }}
{{ $example.Desc }}

```sh 
{{ if $example.WebCLICommand }}command="{{$example.WebCLICommand}}"{{ end }}
{{ $example.Code }}
```
{{ end }}
{{ end }}

{{ range $flagKey, $flagSlice := $subCommand.Flags }}
{{ if $flagSlice }}
### Flags 
{{ range $flag := $flagSlice }}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}{{ end }}{{ end }}{{ end}}{{ end }}
{{ else }}## {{ .Name }}

### Usage

`{{ .Usage }}`
{{ $examples := getExamples . }}{{ if $examples }}
### Examples
{{ range $example := $examples }}
{{ $example.Desc }}

```sh 
{{ if $example.WebCLICommand }}command="{{$example.WebCLICommand}}"{{ end }}
{{ $example.Code }}
```
{{ end }}{{ end }}
{{ range $flagKey, $flagSlice := .Flags }}
### {{ $flagKey }}
{{ range $flag := $flagSlice }}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}{{ end }}{{ end }}