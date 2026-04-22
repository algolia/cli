---
navigation: "cli"
title: |-
  {{ .Name }}
description: |-
  {{ .Description }}
slug: {{ .Slug }}
---
{{ if .Description }}
{{ .Description }}
{{ end }}

## Usage

`{{ .Usage }}`

{{ if .Aliases }}
## Aliases

{{ range $alias := .Aliases -}}
- `{{ $alias }}`
{{ end }}
{{ end }}

{{ if .SubPages }}
## Subcommands

{{ range $subPage := .SubPages -}}
- [`{{ $subPage.Name }}`](/{{ $subPage.Slug }}): {{ $subPage.Description }}
{{ end }}
{{ end }}

{{ $examples := getExamples .Command }}{{ if $examples }}
## Examples
{{ range $example := $examples }}
{{ $example.Desc }}

```sh{{ if $example.WebCLICommand }} command="{{$example.WebCLICommand}}"{{ end }}
{{ $example.Code }}
```
{{ end }}
{{ end }}

{{ range $flagKey, $flagSlice := .Flags -}}
{{ if $flagSlice }}
## {{ $flagKey }}

{{ range $flag := $flagSlice -}}
- {{ if $flag.Shorthand }}`-{{ $flag.Shorthand }}`, {{ end }}`--{{ $flag.Name }}`: {{ $flag.Description }}
{{ end }}
{{ end }}
{{ end }}