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
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument type %s for len not supported",arg.Type())
			}
		},
	},
	"first" :&object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1{
				return newError("wrong number of arguments for first function")
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("first need argument array")
			}
			array := args[0].(*object.Array)
			if len(array.Elements) > 0{
				return array.Elements[0]
			}
			return NULL
		},
	},
	"last" :&object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1{
				return newError("wrong number of arguments for last function")
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("first need argument array")
			}
			array := args[0].(*object.Array)
			if len(array.Elements) > 0{
				return array.Elements[len(array.Elements)-1]
			}
			return NULL
		},
	},
	"tail" :&object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1{
				return newError("wrong number of arguments for tail function")
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("first need argument array")
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if len(array.Elements) > 0{
				newElements := make([]object.Object,length-1,length-1)
				copy(newElements,array.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"head" :&object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1{
				return newError("wrong number of arguments for tail function")
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("first need argument array")
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if len(array.Elements) > 0{
				newElements := make([]object.Object,length-1,length-1)
				copy(newElements,array.Elements[:length-1])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push" :&object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2{
				return newError("wrong number of arguments for push function")
			}
			if args[0].Type() != object.ARRAY_OBJ{
				return newError("first need argument array")
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if len(array.Elements) > 0{
				newElements := make([]object.Object,length+1,length+1)
				copy(newElements,array.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
}