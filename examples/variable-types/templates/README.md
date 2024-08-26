# String variable: {{ .Variables.STRING_VAR }}

# Boolean variable: {{ .Variables.BOOLEAN_VAR }}

# Select variable: {{ .Variables.SELECT_VAR }}

# Multi-select variable: {{ .Variables.MULTI_SELECT_VAR | join ", " }}

# Optional multi-select variable: {{ .Variables.OPTIONAL_MULTI_SELECT_VAR | join ", " }}

# Table variable

| COLUMN_1 | COLUMN_2 | COLUMN_3 |
| --- | --- | --- |
{{- range $val := .Variables.TABLE_VAR }}
| {{ $val.COLUMN_1 }} | {{ $val.COLUMN_2 }} | {{ $val.COLUMN_3 }} |
{{- end}}
