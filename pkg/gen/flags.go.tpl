// This file is generated; DO NOT EDIT.

package cmdutil

import (
    "github.com/spf13/cobra"
    "github.com/MakeNowJust/heredoc"
)

var SearchParams = []string{ {{ range $resName, $resData := .SpecFlags }}
    {{ range $flagName, $flag := $resData.Flags }}"{{ $flagName }}",
    {{ end }}{{ end }}
}

func AddSearchFlags(cmd *cobra.Command) { {{ range $resName, $resData := .SpecFlags }}{{ range $flagName, $flag := $resData.Flags }}{{ if eq $flag.Type "string" }}
        cmd.Flags().String("{{ $flagName }}", {{ if $flag.Def }}"{{ $flag.Def }}"{{ else }}""{{ end }}, heredoc.Doc(`{{ $flag.Usage }}`)){{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ else if eq $flag.Type "boolean" }}
        cmd.Flags().Bool("{{ $flagName }}", {{ $flag.Def }}, heredoc.Doc(`{{ $flag.Usage }}`)){{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ else if eq $flag.Type "integer" }}
        cmd.Flags().Int("{{ $flagName }}", {{ if $flag.Def }}{{ $flag.Def }}{{ else }}0{{ end }}, heredoc.Doc(`{{ $flag.Usage }}`)){{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ else if eq $flag.Type "number" }}
        cmd.Flags().Float64("{{ $flagName }}", {{ if $flag.Def }}{{ $flag.Def }}{{ else }}0{{ end }}, heredoc.Doc(`{{ $flag.Usage }}`)){{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ else if eq $flag.Type "array" }}{{ if eq $flag.SubType "string" }}
        cmd.Flags().StringSlice("{{ $flagName }}", []string{ {{ range $val := $flag.Def }}"{{ $val }}",{{ end }} }, heredoc.Doc(`{{ $flag.Usage }}`)){{ end }}{{ if eq $flag.SubType "integer" }}
        cmd.Flags().IntSlice("{{ $flagName }}", []int{ {{ range $val := $flag.Def }}{{ $val }},{{ end }} }, heredoc.Doc(`{{ $flag.Usage }}`)){{ end }}{{ if eq $flag.SubType "number" }}
        cmd.Flags().Float64Slice("{{ $flagName }}", []float64{ {{ range $val := $flag.Def }}{{ $val }},{{ end }} }, heredoc.Doc(`{{ $flag.Usage }}`)){{ end }}{{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ else }}
        {{ $flagName }} := NewJSONVar([]string{ {{ range $val := $flag.OneOf }}"{{ $val }}",{{ end }} }...)
        cmd.Flags().Var({{ $flagName }}, "{{ $flagName }}", heredoc.Doc(`{{ $flag.Usage }}`)){{ if $flag.Categories }}
        cmd.Flags().SetAnnotation("{{ $flagName }}", "Categories", []string{ {{ range $category := $flag.Categories }}"{{ $category }}", {{ end }} }){{ end }}{{ end }}{{ end }}{{ end }}
}

