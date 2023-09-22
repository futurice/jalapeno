FOO: {{ .Variables.FOO }}
BAR: {{ default "BAR was not set" .Variables.BAR }}
