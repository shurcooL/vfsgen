package bindata

import (
	"regexp"
	"testing"

	"golang.org/x/tools/godoc/vfs"
)

func TestSafeFunctionName(t *testing.T) {
	var knownFuncs = make(map[string]int)
	name1 := safeFunctionName("foo/bar", knownFuncs)
	name2 := safeFunctionName("foo_bar", knownFuncs)
	if name1 == name2 {
		t.Errorf("name collision")
	}
}

func TestFindFiles(t *testing.T) {
	var toc []Asset
	var knownFuncs = make(map[string]int)
	err := findFiles(vfs.OS("./testdata/dupname/"), &toc, []*regexp.Regexp{}, knownFuncs)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}
	if toc[0].Func == toc[1].Func {
		t.Errorf("name collision")
	}
}
