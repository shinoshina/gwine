package evaluator

import (
	"gwine/ast"
	"gwine/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		//fmt.Println("你爹在这")
		//fmt.Println(node.String())
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInflixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}
// 将program 和 blockstatement 分开是因为 如果有嵌套block ，每层block 都return ，由于return 语句返回的只是未封装的值
// 所以无法持续跟踪 returnvalue，这会导致外层 block 将内层block 的return当作一般值来看，导致无法提前退出，最后的结果是外层return的值
// 但其实只要把eval program的返回值重新改为 returnValue 而不是returnValue.Value就行了 ，不这样做的目的是为了让主函数的return和
// block中的return 区分开来，意思是主函数中的返回结果必须是unwrapped return value 而block 只要最后给主函数的是个wrapped return value交给最后流程
// 来unwrap就行了
func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		if returnValue ,ok := result.(*object.ReturnValue);ok{
			return returnValue.Value
		}
	}

	return result
}
func evalBlockStatement(bs *ast.BlockStatement) object.Object{
	var result object.Object

	for _, stmt := range bs.Statements {
		result = Eval(stmt)

		// if returnValue ,ok := result.(*object.ReturnValue);ok{
		// 	return returnValue
		// }
		// 无需解包
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ{
			return result
		}
	}

	return result

}
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}
func evalInflixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInflixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTrue(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}
func evalIntegerInflixExpression(operator string, left, right object.Object) object.Object {
	leftValue, rightValue := left.(*object.Integer).Value, right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}

}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}
func isTrue(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
