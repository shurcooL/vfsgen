// vfsgendev is a convenience tool for using vfsgen in a common development configuration.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	sourceTags = "dev"
	outputTags = "!dev"
)

var (
	sourceFlag = flag.String("source", "", "Specifies the http.FileSystem variable to use as source.")
	nFlag      = flag.Bool("n", false, "Print the generated source but do not run it.")
)

func gen() error {
	source, err := parseSourceFlag(*sourceFlag)
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
