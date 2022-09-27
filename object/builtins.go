package object

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments for len function")
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument type %s for len not supported", arg.Type())
				}
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1{
					return newError("wrong number of arguments for first function")
				}
				if args[0].Type() != ARRAY_OBJ{
					return newError("first need argument array")
				}
				array := args[0].(*Array)
				if len(array.Elements) > 0{
					return array.Elements[0]
				}
				return NullObj
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1{
					return newError("wrong number of arguments for last function")
				}
				if args[0].Type() != ARRAY_OBJ{
					return newError("first need argument array")
				}
				array := args[0].(*Array)
				if len(array.Elements) > 0{
					return array.Elements[len(array.Elements)-1]
				}
				return NullObj
			},
		},
	},
	{
		Name: "head",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1{
					return newError("wrong number of arguments for tail function")
				}
				if args[0].Type() != ARRAY_OBJ{
					return newError("first need argument array")
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if len(array.Elements) > 0{
					newElements := make([]Object,length-1,length-1)
					copy(newElements,array.Elements[:length-1])
					return &Array{Elements: newElements}
				}
				return NullObj
			},
		},
	},
	{
		Name: "push",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2{
					return newError("wrong number of arguments for push function")
				}
				if args[0].Type() != ARRAY_OBJ{
					return newError("first need argument array")
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if len(array.Elements) > 0{
					newElements := make([]Object,length+1,length+1)
					copy(newElements,array.Elements)
					newElements[length] = args[1]
					return &Array{Elements: newElements}
				}
				return NullObj
			},
		},
	},
	{
		Name: "tail",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1{
					return newError("wrong number of arguments for tail function")
				}
				if args[0].Type() != ARRAY_OBJ{
					return newError("first need argument array")
				}
				array := args[0].(*Array)
				length := len(array.Elements)
				if len(array.Elements) > 0{
					newElements := make([]Object,length-1,length-1)
					copy(newElements,array.Elements[1:length])
					return &Array{Elements: newElements}
				}
				return NullObj
			},
		},
	},
}
