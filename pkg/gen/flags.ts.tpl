import { Flags } from "../../../Terminal/types";

const flags: Flags = { {{ range $flagName, $flag := .Flags }}
  {{ $flagName }}: {
    type: {{ if or (eq $flag.Type "string") (eq $flag.SubType "string")}}"string"{{ else if or (eq $flag.Type "boolean") (eq $flag.SubType "boolean")}}"boolean"{{ else if or (eq $flag.Type "integer") (eq $flag.SubType "integer") }}"number"{{ else }}"string"{{ end }},
    shortDesc: `{{ $flag.Usage }}`,
    multiple: {{ if eq $flag.Type "array"}}true{{ else }}false{{ end }},
  },{{ end }}
};

export default flags;
