package main

import (
	"strconv"
	"text/template"
)

type data struct {
	source
	BuildTags string
}

var generateTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"quote": func(s string) string {
		return strconv.Quote(s)
	},
}).Parse(`package main

import (
	"log"

	"github.com/shurcooL/vfsgen"

	{{.ImportPath | quote}}
)

func main() {
	err := vfsgen.Generate({{.PackageName}}.{{.VariableName}}, vfsgen.Options{
		PackageName:     {{.PackageName | quote}},
		BuildTags:       {{.BuildTags | quote}},
		VariableName:    {{.VariableName | quote}},
		VariableComment: {{.VariableComment | quote}},
	})
	if err != nil {
		log.Fatalln(err)
	}
}
`))
