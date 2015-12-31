// vfsgendev is a convenience tool for using vfsgen in a common development configuration.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

var (
	sourceFlag = flag.String("source", "", "Specifies the http.FileSystem variable to use as source.")
	nFlag      = flag.Bool("n", false, "Print the generated source but do not run it.")
)

const (
	sourceTags = "dev"

	outputFilename = "data_generate.go"
	outputTags     = "!dev"
)

type source struct {
	ImportPath   string
	PackageName  string
	VariableName string
}

var errInvalidFormat = fmt.Errorf("invalid format")

// parseSourceFlag parses the "-source" flag. It must have "import/path".VariableName format.
func parseSourceFlag(sourceFlag string) (source, error) {
	// Parse sourceFlag as a Go expression, albeit a strange one:
	//
	// 	"import/path".VariableName
	//
	e, err := parser.ParseExpr(sourceFlag)
	if err != nil {
		return source{}, errInvalidFormat
	}
	se, ok := e.(*ast.SelectorExpr)
	if !ok {
		return source{}, errInvalidFormat
	}
	importPath, err := stringValue(se.X)
	if err != nil {
		return source{}, err
	}
	variableName := se.Sel.Name

	// Import package to get its full import path and package name.
	ctx := build.Default
	ctx.BuildTags = strings.Fields(sourceTags)
	bpkg, err := ctx.Import(importPath, ".", 0)
	if err != nil {
		return source{}, fmt.Errorf("can't import package %q: %v", importPath, err)
	}

	return source{
		ImportPath:   bpkg.ImportPath,
		PackageName:  bpkg.Name,
		VariableName: variableName,
	}, nil
}

// stringValue returns the string value of string literal e.
func stringValue(e ast.Expr) (string, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		return "", errInvalidFormat
	}
	if lit.Kind != token.STRING {
		return "", errInvalidFormat
	}
	return strconv.Unquote(lit.Value)
}

var t = template.Must(template.New("").Funcs(template.FuncMap{
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
		PackageName:  {{.PackageName | quote}},{{with .BuildTags}}
		BuildTags:    {{. | quote}},{{end}}
		VariableName: {{.VariableName | quote}},
	})
	if err != nil {
		log.Fatalln(err)
	}
}
`))

type data struct {
	source
}

func (data) BuildTags() string { return outputTags }

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

func gen() error {
	source, err := parseSourceFlag(*sourceFlag)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data{source: source})
	if err != nil {
		return err
	}

	if *nFlag == true {
		io.Copy(os.Stdout, &buf)
		return nil
	}

	err = run(buf.String(), sourceTags)
	return err
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: vfsgendev [flags] -source="import/path".VariableName`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
		return
	}

	err := gen()
	if err != nil {
		log.Fatalln(err)
	}
}
