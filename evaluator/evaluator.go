package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

var (
	NULL  = object.NIL
	TRUE  = object.TRUE
	FALSE = object.FALSE
)

// Evaluator is a structure that will define methods to be used
// to evaluate the AST nodes.
type Evaluator struct {
	pos      token.Pos
	err      *object.Error
	errStack error
}

// Eval walks through the AST and evaluates the nodes into an object
func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	if node == nil {
		return nil
	}
	e.pos = node.Pos()
	switch node := node.(type) {
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return object.NewBoolean(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfStmt:
		return e.evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)
		right := e.Eval(node.Right, env)
		return e.evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left := e.Eval(node.Left, env)
		result := e.evalPostfixExpression(node.Operator, left)
		if _, ok := node.Left.(*ast.Identifier); ok && !e.isError(result) { // update identifier
			env.Update(node.Left.Literal(), result)
		}
		return result
	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)
		return e.evalPrefixExpression(node.Operator, right)
	case *ast.ReturnExpression:
		return e.evalReturnExpression(node, env)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.LetStmt:
		if expr, ok := node.Expr.(*ast.MatchExpression); ok {
			return e.evalMatchExpression(&node.Name.Value, expr, env)
		}
		return e.evalLetExpression(node.Name.Value, node.Expr, env)
	case *ast.MatchExpression:
		return e.evalMatchExpression(nil, node, env)
	case *ast.BlockStmt:
		blockEnv := object.NewEnvironment(env)
		return e.evalBlockStmt(node, blockEnv)
	case *ast.ExpressionStmt:
		return e.Eval(node.Expr, env)
	case *ast.ConditionalStmt:
		return e.Eval(node.Statement, env)
	case *ast.ImportStmt:
		return e.evalImportStmt(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Body: node.Body, Params: node.Params, ParentEnv: env}
	case *ast.CallExpression:
		fn := e.Eval(node.Function, env)
		if e.isError(fn) {
			return fn
		}
		args := e.evalArgs(node.Args, env)
		return e.applyFunction(fn, args)
	case *ast.ArrayLiteral:
		entries := e.evalArgs(node.Elements, env)
		return &object.Array{Entries: &entries}
	case *ast.RangeArrayLiteral:
		return e.evalRangeArray(node, env)
	case *ast.HashLiteral:
		entries := e.evalPairs(node.Pair, env)
		for _, val := range entries {
			if e.isError(val) {
				return val
			}
		}
		return &object.Hash{Entries: entries}
	case *ast.SetLiteral:
		entries := e.evalSet(node.Elements, env)
		for key := range entries {
			if key == object.EmptyHashKey {
				return e.err
			}
		}
		return &object.Set{Entries: entries}
	case *ast.ReassignmentStmt:
		return e.evalReassignmentStmt(node, env)
	case *ast.ForLoopStmt:
		return e.evalForLoopStmt(node, env)
	case *ast.ObjectMethodExpression:
		return e.evalObjectDotExpr(node, env)
	}

	return nil
}

func (e *Evaluator) evalReturnExpression(node *ast.ReturnExpression, env *object.Environment) object.Object {
	returnVal := e.Eval(node.Expr, env)
	if returnVal.Type() == object.ERROR_OBJ {
		return returnVal
	}
	// wrap the value so that block statements can terminate early if they encounter a return
	return &object.ReturnValue{Value: returnVal}
}

// evalArgs evaluates arguments
func (e *Evaluator) evalArgs(args []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(args))

	for _, arg := range args {
		result = append(result, e.Eval(arg, env))
	}

	return result
}

func (e *Evaluator) evalPairs(args map[ast.Expression]ast.Expression, env *object.Environment) map[string]object.Object {
	result := make(map[string]object.Object)

	for key, value := range args {
		keyObj := e.Eval(key, env)
		keyString := object.ToRawValue(keyObj)
		if keyString == "" {
			result[keyString] = object.NewErrorWithMsg(fmt.Sprintf("invalid key '%s'", keyString))
			return result
		}
		valueObj := e.Eval(value, env)
		if e.isError(keyObj) || e.isError(valueObj) {
			result[keyString] = valueObj
			return result
		}
		result[keyString] = valueObj
	}

	return result
}

func (e *Evaluator) evalSet(args map[ast.Expression]struct{}, env *object.Environment) map[object.HashKey]struct{} {
	result := make(map[object.HashKey]struct{})

	for key := range args {
		keyObj := e.Eval(key, env)
		hashKey := object.ToHashKey(keyObj)
		if hashKey == object.EmptyHashKey {
			e.err = object.NewErrorWithMsg(fmt.Sprintf("invalid set entry '%s'", keyObj.Inspect()))
			keyObj = e.err
		}

		if e.isError(e.err) {
			result[object.EmptyHashKey] = struct{}{}
			return result
		}
		result[hashKey] = struct{}{}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	if b, ok := builtins[node.Value]; ok {
		return b
	}

	return nil
}

// evalBlockStmt evaluates the block statement. The env passed here
// should be one scoped to the block only
func (e *Evaluator) evalBlockStmt(node *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}
	for _, stmt := range node.Statements {
		result = e.Eval(stmt, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func (e *Evaluator) evalImportStmt(node *ast.ImportStmt, env *object.Environment) object.Object {
	if node == nil {
		return object.NewErrorWithMsg("invalid import") //TODO improve error message
	}

	if mod, ok := modules[node.Value]; ok {
		env.Set(node.Value, object.NewImport(mod))
		return NULL
	}
	return object.NewErrorWithMsg("invalid import") //TODO improve error message
}

// Indexable defines an interface for objects that can be indexed
type Indexable interface {
	Update(object.Object, object.Object) object.Object
}

func (e *Evaluator) evalReassignmentStmt(node *ast.ReassignmentStmt, env *object.Environment) object.Object {
	switch expr := node.Name.(type) {
	case *ast.Identifier:
		if _, found := env.Get(expr.Value); !found {
			return e.EvalError(fmt.Sprintf("cannot reassign undeclared identifier '%s'", expr.Value), expr.Pos())
		}
		res := e.Eval(node.Expr, env)
		env.Update(expr.Value, res)
		return res
	case *ast.IndexExpression:
		// evaluate the left, the index, and then set the left[index] to the RHS of the reassignment
		left := e.Eval(expr.Left, env)
		if e.isError(left) {
			return left
		}
		leftIndexable, ok := left.(Indexable)
		if !ok {
			return object.NewErrorWithMsg("object of type %T not indexable", left)
		}
		index := e.Eval(expr.Index, env)
		if e.isError(index) {
			return index
		}
		rhs := e.Eval(node.Expr, env)
		if e.isError(rhs) {
			return rhs
		}

		resp := leftIndexable.Update(index, rhs)
		if e.isError(resp) {
			return resp
		}
	default:
		return object.NewErrorWithMsg("invalid reassignment")
	}
	return NULL
}

func (e *Evaluator) evalIfExpression(node *ast.IfStmt, env *object.Environment) object.Object {
	cond := e.Eval(node.Consequence.Condition, env)
	if isTruthy(cond) {
		return e.Eval(node.Consequence, env)
	} else {
		for _, alt := range node.Alternatives {
			if alt.Condition == nil { // normal else branch (else)
				return e.Eval(alt, env)
			}
			cond := e.Eval(alt.Condition, env) // (else if branch)
			if isTruthy(cond) {
				return e.Eval(alt, env)
			}
		}
	}
	return NULL
}

func (e *Evaluator) evalLetExpression(nodeName string, RHS ast.Expression, env *object.Environment) object.Object {
	expr := e.Eval(RHS, env)
	if !e.isError(expr) {
		env.Set(nodeName, expr)
		return NULL
	}
	// if the let expr is an error, we return it so we can terminate the program
	return expr
}

func (e *Evaluator) evalMatchExpression(exprIdent *string, node *ast.MatchExpression, env *object.Environment) object.Object {
	// evaluate the match expression i.e match(expr_here)
	expr := e.Eval(node.Expression, env)
	if expr == nil {
		return object.NewErrorWithMsg("match expression cannot be null")
	}

	// create an environment for the match block, and
	// set the error value for the block
	matchEnv := object.NewEnvironment(env)
	if e.isError(expr) {
		matchEnv.Set(token.ErrorIdentifier, expr)
	} else if exprIdent != nil {
		// if the match is called from a let statement,
		// and no error, set its expression to the identifier
		env.Set(*exprIdent, expr)
	}

	for _, matchCase := range node.Cases {
		// evaluate each case
		pattern := e.Eval(matchCase.Pattern, matchEnv)

		// if the case matches the match expression, return the case output
		if pattern != nil && pattern.Equal(expr) {
			val := e.Eval(matchCase.Output, matchEnv)
			if _, found := matchEnv.Get(token.ErrorIdentifier); found {
				matchEnv.Set(token.ErrorIdentifier, object.NewString(fmt.Sprintf("error: %s", expr.Inspect())))
			}
			return val
		}
		// it is important to check this after the equality check, so we
		// can differentiate other runtime errors from the one returned from the match expression
		if e.isError(pattern) {
			return pattern
		}
	}

	// if no case matches and there is a default block
	if node.Default != nil {
		return e.Eval(node.Default, matchEnv)
	}

	// if no case matches, set the expr to the let stmt
	env.Set(*exprIdent, expr)
	return NULL
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := e.Eval(node.Left, env)
	if e.isError(left) {
		return left
	}

	switch left.Type() {
	case object.ARRAY_OBJ:
		left := left.(*object.Array)
		if index, ok := node.Index.(*ast.IntegerLiteral); ok {
			if int(index.Value) >= len(*left.Entries) {
				return e.EvalError(fmt.Sprintf("index %d out of range with length %d", index.Value, len(*left.Entries)), node.Pos())
			}
			return (*left.Entries)[index.Value]
		}
	case object.HASH_OBJ:
		left := left.(*object.Hash)
		index := node.Index.(*ast.StringLiteral)
		if entry, ok := left.Entries[index.Value]; ok {
			return entry
		}
	}
	return e.EvalError(fmt.Sprintf("invalid index entry '%s' for '%s'", node.Index.Literal(), node.Left.Literal()), node.Pos())
}

func (e *Evaluator) evalPostfixExpression(operator string, left object.Object) object.Object {
	if e.isError(left) {
		return left
	}
	if left.Type() == object.INT_OBJ {
		left := left.(*object.Int)
		switch operator {
		case token.INC:
			return &object.Int{Value: left.Value + 1}
		case token.DEC:
			return &object.Int{Value: left.Value - 1}
		}
	}
	if left.Type() == object.FLOAT_OBJ {
		left := left.(*object.Float)
		switch operator {
		case token.INC:
			return &object.Float{Value: left.Value + 1}
		case token.DEC:
			return &object.Float{Value: left.Value - 1}
		}
	}
	return e.EvalError(fmt.Sprintf("invalid postfix operator %s for %s", operator, left.Inspect()), e.pos)
}

func (e *Evaluator) evalBangOperator(operator string, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Int:
		return e.booleanObj(right.Value == 0)
	case *object.String:
		return e.booleanObj(len(right.Value) == 0)
	case *object.Array:
		return e.booleanObj(len(*right.Entries) == 0)
	case *object.Boolean:
		return e.booleanObj(!right.Value)
	case *object.Nil:
		return TRUE
	}
	return e.EvalError(fmt.Sprintf("bang operator not valid for type %T", right), e.pos)
}

func (e *Evaluator) booleanObj(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) object.Object {
	if e.isError(fn) {
		return fn
	}
	// for _, arg := range args {
	// 	if e.isError(arg) {
	// 		return arg
	// 	}
	// }
	switch fn := fn.(type) {
	case *object.Function:
		fnEnv := object.NewEnvironment(fn.ParentEnv)
		for i, p := range fn.Params {
			fnEnv.Set(p.Value, args[i])
		}
		result := e.Eval(fn.Body, fnEnv)
		return unwrapReturnValue(result)
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj.Type() {
	case object.BOOLEAN_OBJ:
		obj := obj.(*object.Boolean)
		return obj.Value
	case object.INT_OBJ:
		obj := obj.(*object.Int)
		return obj.Value != 0
	case object.ERROR_OBJ:
		obj := obj.(*object.Error)
		return obj.Message == ""
	}
	return true
}

// isError returns true if the object is an error. If the object is nil, error is false
func (e *Evaluator) isError(obj object.Object) bool {
	if obj != nil {
		errObj, ok := obj.(*object.Error)
		if !ok {
			return false
		}
		if errObj == nil {
			return false
		}
		e.errStack = multierror.Append(e.errStack, errObj.Native().(error))
		return true
	}
	return false
}

func (e *Evaluator) EvalError(err string, pos token.Pos) *object.Error {
	msg := fmt.Sprintf(`
	Error: %s
	Line: %d
	Column: %d
	`, err, pos.Line, pos.Column)
	return &object.Error{Message: msg}
}
