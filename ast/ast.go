// nolint: go-staticcheck
package ast

import (
	"ede/token"
	"fmt"
)

type Node interface {
	Pos() token.Pos
	Literal() string
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type (
	LetStmt struct {
		ValuePos token.Pos
		Name     *Identifier
		Expr     Expression
		Token    token.Token
	}

	ExpressionStmt struct {
		ValuePos token.Pos
		Expr     Expression
		Token    token.Token
	}

	BlockStmt struct {
		Statements []Statement
		ValuePos   token.Pos
	}
	ConditionalStmt struct {
		Token     token.Token
		Condition Expression
		Statement Statement
		ValuePos  token.Pos
	}

	Identifier struct {
		Token    token.Token
		Value    string
		ValuePos token.Pos
	}

	StringLiteral struct {
		Token    token.Token
		Value    string
		ValuePos token.Pos
	}

	IntegerLiteral struct {
		Token    token.Token
		Value    int64
		ValuePos token.Pos
	}

	BooleanLiteral struct {
		Token    token.Token
		Value    bool
		ValuePos token.Pos
	}

	FunctionLiteral struct {
		Token    token.Token
		Params   []*Identifier
		Body     *BlockStmt
		ValuePos token.Pos
	}
	ArrayLiteral struct {
		Token    token.Token
		Elements []Expression
		ValuePos token.Pos
	}

	IfExpression struct {
		Condition    Expression
		Consequence  *ConditionalStmt
		Alternatives []*ConditionalStmt
		ValuePos     token.Pos
		Token        token.Token
	}
	InfixExpression struct {
		Left     Expression
		Right    Expression
		Operator string
		ValuePos token.Pos
		Token    token.Token
	}

	PrefixExpression struct {
		Operator string
		Token    token.Token
		Right    Expression
		ValuePos token.Pos
	}

	PostfixExpression struct {
		Operator string
		Token    token.Token
		Left     Expression
		ValuePos token.Pos
	}

	ReturnExpression struct {
		ValuePos token.Pos
		Expr     Expression
		Token    token.Token
	}

	IndexExpression struct {
		ValuePos token.Pos
		Left     Expression
		Index    Expression
		Token    token.Token
	}

	CallExpression struct {
		Function Expression // this can be function literal or identifier
		Args     []Expression
		Token    token.Token
		ValuePos token.Pos
	}

	Program struct {
		ParseErrors []error
		Statements  []Statement
		ValuePos    token.Pos
	}
)

func (s *LetStmt) stmtNode()           {}
func (s *ExpressionStmt) stmtNode()    {}
func (s *BlockStmt) stmtNode()         {}
func (s *ConditionalStmt) stmtNode()   {}
func (s *StringLiteral) stmtNode()     {}
func (s *IntegerLiteral) stmtNode()    {}
func (s *ArrayLiteral) stmtNode()      {}
func (s *BooleanLiteral) stmtNode()    {}
func (s *Identifier) stmtNode()        {}
func (s *InfixExpression) stmtNode()   {}
func (s *IfExpression) stmtNode()      {}
func (s *PrefixExpression) stmtNode()  {}
func (s *ReturnExpression) stmtNode()  {}
func (s *PostfixExpression) stmtNode() {}
func (s *CallExpression) stmtNode()    {}
func (s *IndexExpression) stmtNode()   {}

func (s *StringLiteral) exprNode()     {}
func (s *FunctionLiteral) exprNode()   {}
func (s *IntegerLiteral) exprNode()    {}
func (s *ArrayLiteral) exprNode()      {}
func (s *BooleanLiteral) exprNode()    {}
func (s *Identifier) exprNode()        {}
func (s *InfixExpression) exprNode()   {}
func (s *IfExpression) exprNode()      {}
func (s *PrefixExpression) exprNode()  {}
func (s *ReturnExpression) exprNode()  {}
func (s *PostfixExpression) exprNode() {}
func (s *CallExpression) exprNode()    {}
func (s *IndexExpression) exprNode()   {}

func (s *Program) Pos() token.Pos           { return s.ValuePos }
func (s *LetStmt) Pos() token.Pos           { return s.ValuePos }
func (s *ExpressionStmt) Pos() token.Pos    { return s.ValuePos }
func (s *BlockStmt) Pos() token.Pos         { return s.ValuePos }
func (s *ConditionalStmt) Pos() token.Pos   { return s.ValuePos }
func (s *StringLiteral) Pos() token.Pos     { return s.ValuePos }
func (s *FunctionLiteral) Pos() token.Pos   { return s.ValuePos }
func (s *IntegerLiteral) Pos() token.Pos    { return s.ValuePos }
func (s *ArrayLiteral) Pos() token.Pos      { return s.ValuePos }
func (s *BooleanLiteral) Pos() token.Pos    { return s.ValuePos }
func (s *Identifier) Pos() token.Pos        { return s.ValuePos }
func (s *InfixExpression) Pos() token.Pos   { return s.ValuePos }
func (s *IfExpression) Pos() token.Pos      { return s.ValuePos }
func (s *PrefixExpression) Pos() token.Pos  { return s.ValuePos }
func (s *ReturnExpression) Pos() token.Pos  { return s.ValuePos }
func (s *PostfixExpression) Pos() token.Pos { return s.ValuePos }
func (s *CallExpression) Pos() token.Pos    { return s.ValuePos }
func (s *IndexExpression) Pos() token.Pos   { return s.ValuePos }

func (s *Program) Literal() string         { return "" } // TODO
func (s *LetStmt) Literal() string         { return s.Token.Literal }
func (s *ExpressionStmt) Literal() string  { return s.Token.Literal }
func (s *BlockStmt) Literal() string       { return "" } // TODO
func (s *ConditionalStmt) Literal() string { return "" } // TODO
func (s *StringLiteral) Literal() string   { return s.Value }
func (s *FunctionLiteral) Literal() string { return s.Token.Literal } //TODO
func (s *IntegerLiteral) Literal() string  { return fmt.Sprint(s.Value) }
func (s *ArrayLiteral) Literal() string    { return "" } // TODO
func (s *BooleanLiteral) Literal() string  { return fmt.Sprint(s.Value) }
func (s *Identifier) Literal() string      { return s.Value }
func (s *InfixExpression) Literal() string {
	return fmt.Sprintf("(%s %s %s)", s.Left.Literal(), s.Operator, s.Right.Literal())
}
func (s *IfExpression) Literal() string { return s.Token.Literal }
func (s *PrefixExpression) Literal() string {
	return fmt.Sprintf("%s%s", s.Token.Literal, s.Right.Literal())
}
func (s *ReturnExpression) Literal() string {
	return fmt.Sprintf("%s%s", s.Token.Literal, s.Expr.Literal())
}
func (s *PostfixExpression) Literal() string { return s.Token.Literal }
func (s *CallExpression) Literal() string    { return s.Token.Literal }
func (s *IndexExpression) Literal() string   { return s.Token.Literal }
