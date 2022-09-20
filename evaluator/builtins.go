package evaluator

import(
	"gwine/object"
)
var builtins = map[string]*object.Builtin{
	"len" : &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1{
				return newError("wrong number of arguments for len function")
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument type %s for len not supported",arg.Type())
			}
		},
	},
}