package config

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// isMaxBinding — returns true if the first argument is greater than all others.
func isMaxBinding(args ...ref.Val) ref.Val {
	if len(args) < 2 {
		return types.Bool(false)
	}
	first, ok := args[0].Value().(int64)
	if !ok {
		return types.Bool(false)
	}
	for _, v := range args[1:] {
		val, ok := v.Value().(int64)
		if !ok || first <= val {
			return types.Bool(false)
		}
	}
	return types.Bool(true)
}

// IsMaxLib — EnvOption that registers the isMax function for CEL.
var IsMaxLib = cel.Function(
	"isMax",
	cel.Overload(
		"isMax_int_varargs_bool",
		[]*cel.Type{cel.IntType, cel.DynType},
		cel.BoolType,
		cel.FunctionBinding(func(args ...ref.Val) ref.Val {
			return isMaxBinding(args...)
		}),
	),
)
