{{template "header"}}
import flixTemplate from "{{.Location}}"

const parameterNames = [
  {{- range $index, $ele := .Parameters -}}
     {{if $index}}, {{end}}"{{$ele.Name}}"
  {{- end -}}
];

// TODO: test parameters haven't changed at runtime
export async function {{.Title}}({ 
 {{- if len .Parameters -}}
    {{- range $index, $ele := .Parameters -}}
      {{if $index}}, {{end}}{{.Name}}
    {{- end -}}
  {{- end -}}
}) {
  const info = await fcl.query({
    template: flixTemplate,
    {{ if len .Parameters -}}
    args: (arg, t) => [
      {{- range $index, $ele := .Parameters -}}
        {{if $index}}, {{end}}arg({{.Name}}, t.{{.Type}})
      {{- end -}}
      ]
    {{- end }}
  });

  return info
}
