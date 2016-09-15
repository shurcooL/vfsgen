package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type source struct {
	ImportPath      string
	PackageName     string
	VariableName    string
	VariableComment string
}

// parseSourceFlag parses the "-source" flag. It must have "import/path".VariableName format.
func parseSourceFlag(sourceFlag string) (source, error) {
	// Parse sourceFlag as a Go expression, albeit a strange one:
	//
	// 	"import/path".VariableName
	//
	e, err := parser.ParseExpr(sourceFlag)
	if err != nil {
		return source{}, fmt.Errorf("invalid format")
	}
	se, ok := e.(*ast.SelectorExpr)
	if !ok {
		return source{}, fmt.Errorf("invalid format")
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
	dpkg, err := docPackage(bpkg)
	if err != nil {
		return source{}, fmt.Errorf("can't get godoc of package %q: %v", importPath, err)
	}
	var variableComment string
	for _, v := range dpkg.Vars {
		if len(v.Names) == 1 && v.Names[0] == variableName {
			variableComment = strings.TrimSuffix(v.Doc, "\n")
			break
		}
	}

	return source{
		ImportPath:      bpkg.ImportPath,
		PackageName:     bpkg.Name,
		VariableName:    variableName,
		VariableComment: variableComment,
	}, nil
}

// stringValue returns the string value of string literal e.
func stringValue(e ast.Expr) (string, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		return "", fmt.Errorf("invalid format")
	}
	if lit.Kind != token.STRING {
		return "", fmt.Errorf("invalid format")
	}
	return strconv.Unquote(lit.Value)
}
