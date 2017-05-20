package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

// parseSourceFlag parses the "-source" flag value. It must have "import/path".VariableName format.
func parseSourceFlag(sourceFlag string) (importPath, variableName string, err error) {
	// Parse sourceFlag as a Go expression, albeit a strange one:
	//
	// 	"import/path".VariableName
	//
	e, err := parser.ParseExpr(sourceFlag)
	if err != nil {
		return "", "", fmt.Errorf("invalid format")
	}
	se, ok := e.(*ast.SelectorExpr)
	if !ok {
		return "", "", fmt.Errorf("invalid format")
	}
	importPath, err = stringValue(se.X)
	if err != nil {
		return "", "", err
	}
	variableName = se.Sel.Name
	return importPath, variableName, nil
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

type source struct {
	ImportPath      string
	PackageName     string
	VariableName    string
	VariableComment string
}

// lookupSource imports package using provided build context, and returns the
// import path, package name, variable name and variable comment.
func lookupSource(bctx build.Context, importPath, variableName string) (source, error) {
	bpkg, err := bctx.Import(importPath, "", 0)
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

// TODO: Keep in sync or unify with github.com/shurcooL/cmd/gorepogen/docpackage.go.
// TODO: See if these can be cleaned up.

func docPackage(bpkg *build.Package) (*doc.Package, error) {
	apkg, err := astPackage(bpkg)
	if err != nil {
		return nil, err
	}
	return doc.New(apkg, bpkg.ImportPath, 0), nil
}

func astPackage(bpkg *build.Package) (*ast.Package, error) {
	// TODO: Either find a way to use golang.org/x/tools/importer (from Go 1.4~ or older, it no longer exists as of Go 1.6) directly, or do file AST parsing in parallel like it does.
	filenames := append(bpkg.GoFiles, bpkg.CgoFiles...)
	files := make(map[string]*ast.File, len(filenames))
	fset := token.NewFileSet()
	for _, filename := range filenames {
		fileAst, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, filename), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[filename] = fileAst // TODO: Figure out if filename or full path are to be used (the key of this map doesn't seem to be used anywhere).
	}
	return &ast.Package{Name: bpkg.Name, Files: files}, nil
}
