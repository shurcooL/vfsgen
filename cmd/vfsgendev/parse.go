package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type source struct {
	ImportPath   string
	PackageName  string
	VariableName string
}

var errInvalidFormat = errors.New("invalid format")

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
