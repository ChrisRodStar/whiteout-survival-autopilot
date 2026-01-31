package config

import (
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// -----------------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------------

// FuzzySubstringMatch returns true if needle occurs in haystack
// with at most maxDist Levenshtein errors within any window of the same length.
// Additionally: if haystack is shorter than needle, compare the strings entirely.
func FuzzySubstringMatch(haystack, needle string, maxDist int) bool {
	haystack = strings.ToLower(haystack)
	needle = strings.ToLower(needle)

	n, m := len(needle), len(haystack)
	if n == 0 {
		return false
	}

	// haystack is shorter than needle → compare entirely
	if m < n {
		return levenshtein.ComputeDistance(haystack, needle) <= maxDist
	}

	for i := 0; i <= m-n; i++ {
		if levenshtein.ComputeDistance(haystack[i:i+n], needle) <= maxDist {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// compareText(a, b) → bool  (registration in CEL)
// -----------------------------------------------------------------------------

// compareTextBinding — actual implementation of the function for CEL.
func compareTextBinding(lhs, rhs ref.Val) ref.Val {
	a, ok1 := lhs.Value().(string)
	b, ok2 := rhs.Value().(string)
	if !ok1 || !ok2 {
		return types.Bool(false)
	}

	al := strings.ToLower(a)
	bl := strings.ToLower(b)

	if bl == "" || al == "" {
		return types.Bool(false)
	}

	if strings.Contains(bl, al) {
		return types.Bool(true)
	}

	// Only for substrings of length >= 4 characters apply fuzzy matching
	if len(al) >= 4 && FuzzySubstringMatch(bl, al, 1) {
		return types.Bool(true)
	}

	return types.Bool(false)
}

// CompareTextLib — EnvOption that registers the compareText function.
var CompareTextLib = cel.Function(
	"compareText",
	cel.Overload(
		"compareText_string_string_bool",
		[]*cel.Type{cel.StringType, cel.StringType},
		cel.BoolType,
		cel.BinaryBinding(compareTextBinding),
	),
)
