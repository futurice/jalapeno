| Environment | Application | Scale |
{{- range $val := .Variables.ENVIRONMENTS}}
| {{$val.name}} | {{$val.app}} | {{$val.scale}} |
{{- end}}
