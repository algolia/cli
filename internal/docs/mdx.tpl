---
title: {{ .Name }}
description: {{ .Description }}
public: true
---

```txt Usage
{{ .Usage }}
```
{{- if .SubPages }}

## Commands

{{ range $subPage := .SubPages -}}
- [`{{ $subPage.Name }}`](/{{ $subPage.Slug }})
{{ end -}}
{{ end -}}
{{ $examples := getExamples .Command -}}
{{ if $examples }}

## Examples

{{ range $example := $examples -}}
{{ if $example.Desc }}{{ $example.Desc }}:

{{ end -}}
```sh icon=square-terminal
{{ $example.Code }}
```
{{ end -}}
{{ end -}}
{{ range $flagKey, $flagSlice := .Flags }}
{{ if $flagSlice }}
## {{ $flagKey }}

{{ range $flag := $flagSlice -}}
<ParamField body="{{ if $flag.Shorthand }}-{{ $flag.Shorthand }}, {{ end }}--{{ $flag.Name }}">

{{ formatAlgoliaDocLinks (trimTrailingNewlines $flag.Description) }}

</ParamField>

{{ end -}}
{{ end }}
{{- end -}}
