package main

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
)

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
