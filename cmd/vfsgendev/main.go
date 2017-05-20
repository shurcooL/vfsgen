// vfsgendev is a convenience tool for using vfsgen in a common development configuration.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	sourceFlag = flag.String("source", "", "Specifies the http.FileSystem variable to use as source.")
	nFlag      = flag.Bool("n", false, "Print the generated source but do not run it.")
)

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
	}

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

const (
	sourceTags = "dev"
	outputTags = "!dev"
)

func run() error {
	importPath, variableName, err := parseSourceFlag(*sourceFlag)
	if err != nil {
		return err
	}

	bctx := build.Default
	bctx.BuildTags = strings.Fields(sourceTags)
	source, err := lookupSource(bctx, importPath, variableName)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = generateTemplate.Execute(&buf, data{source: source})
	if err != nil {
		return err
	}

	if *nFlag {
		io.Copy(os.Stdout, &buf)
		return nil
	}

	err = goRun(buf.String(), sourceTags)
	return err
}

// goRun runs Go code src with build tags.
func goRun(src string, tags string) error {
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