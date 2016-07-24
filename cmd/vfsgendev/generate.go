package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"text/template"
)

type data struct {
	source
}

func (data) BuildTags() string { return outputTags }

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

// run runs Go code src with build tags.
func run(src string, tags string) error {
	// Create a temp folder.
	tempDir, err := ioutil.TempDir("", "vfsgendev_")
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warning: error removing temp dir:", err)
		}
	}()

	// Write the source code file.
	tempFile := filepath.Join(tempDir, "generate.go")
	err = ioutil.WriteFile(tempFile, []byte(src), 0600)
	if err != nil {
		return err
	}

	// Compile and run the program.
	cmd := exec.Command("go", "run", "-tags="+tags, tempFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
