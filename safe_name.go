package bindata

import (
	"fmt"
	"strings"
	"unicode"
)

// safeFunctionName converts the given name into a name
// which qualifies as a valid function identifier. It
// also compares against a known list of functions to
// prevent conflict based on name translation.
func safeFunctionName(name string, knownFuncs map[string]int) string {
	name = strings.ToLower(name)
	name = regFuncName.ReplaceAllString(name, "_")

	// Get rid of "__" instances for niceness.
	for strings.Index(name, "__") > -1 {
		name = strings.Replace(name, "__", "_", -1)
	}

	// Leading underscores are silly (unless they prefix a digit (see below)).
	for len(name) > 1 && name[0] == '_' {
		name = name[1:]
	}

	// Identifier can't start with a digit.
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}

	if num, ok := knownFuncs[name]; ok {
		knownFuncs[name] = num + 1
		name = fmt.Sprintf("%s%d", name, num)
	} else {
		knownFuncs[name] = 2
	}

	return name
}
