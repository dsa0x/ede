// nolint: go-staticcheck
package ast

import (
	"ede/token"
	"fmt"
)

type Node interface {
	Pos() token.Pos
	Literal() string
	TokenType() token.TokenType
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

	ReassignmentStmt struct {
		ValuePos token.Pos
		Name     *Identifier
		Expr     Expression
		Token    token.Token
	}

	BlockStmt struct {
		Token      token.Token
		Statements []Statement
		ValuePos   token.Pos
	}
	ConditionalStmt struct {
		Token     token.Token
		Condition Expression
		Statement Statement
		ValuePos  token.Pos
	}

	ForLoopStmt struct {
		Token     token.Token
		Variable  *Identifier
		Boundary  Expression
		Statement *BlockStmt
		ValuePos  token.Pos
	}

	IfStmt struct {
		Condition    Expression
		Consequence  *ConditionalStmt
		Alternatives []*ConditionalStmt
		ValuePos     token.Pos
		Token        token.Token
	}

	CommentStmt struct {
		Value    string
		ValuePos token.Pos
		Token    token.Token
	}

	ErrorStmt struct {
		Value    string
		ValuePos token.Pos
		Token    token.Token
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

	FloatLiteral struct {
		Token    token.Token
		Value    float64
		ValuePos token.Pos
	}

	BooleanLiteral struct {
		Token    token.Token
		Value    bool
		ValuePos token.Pos
	}

	RangeLiteral struct {
		Token    token.Token
		Start    int64
		End      int64
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

	HashLiteral struct {
		Token    token.Token
		Pair     map[Expression]Expression
		ValuePos token.Pos
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
		Token       token.Token
		ParseErrors []error
		Statements  []Statement
		ValuePos    token.Pos
	}
)

func (s *LetStmt) stmtNode()           {}
func (s *ExpressionStmt) stmtNode()    {}
func (s *BlockStmt) stmtNode()         {}
func (s *CommentStmt) stmtNode()       {}
func (s *ForLoopStmt) stmtNode()       {}
func (s *ConditionalStmt) stmtNode()   {}
func (s *StringLiteral) stmtNode()     {}
func (s *IntegerLiteral) stmtNode()    {}
func (s *ArrayLiteral) stmtNode()      {}
func (s *HashLiteral) stmtNode()       {}
func (s *BooleanLiteral) stmtNode()    {}
func (s *FloatLiteral) stmtNode()      {}
func (s *Identifier) stmtNode()        {}
func (s *ReassignmentStmt) stmtNode()  {}
func (s *ErrorStmt) stmtNode()         {}
func (s *InfixExpression) stmtNode()   {}
func (s *IfStmt) stmtNode()            {}
func (s *PrefixExpression) stmtNode()  {}
func (s *ReturnExpression) stmtNode()  {}
func (s *PostfixExpression) stmtNode() {}
func (s *CallExpression) stmtNode()    {}
func (s *IndexExpression) stmtNode()   {}

func (s *StringLiteral) exprNode()     {}
func (s *FunctionLiteral) exprNode()   {}
func (s *IntegerLiteral) exprNode()    {}
func (s *ArrayLiteral) exprNode()      {}
func (s *HashLiteral) exprNode()       {}
func (s *BooleanLiteral) exprNode()    {}
func (s *FloatLiteral) exprNode()      {}
func (s *Identifier) exprNode()        {}
func (s *ReassignmentStmt) exprNode()  {}
func (s *ErrorStmt) exprNode()         {}
func (s *InfixExpression) exprNode()   {}
func (s *PrefixExpression) exprNode()  {}
func (s *ReturnExpression) exprNode()  {}
func (s *PostfixExpression) exprNode() {}
func (s *CallExpression) exprNode()    {}
func (s *IndexExpression) exprNode()   {}

func (s *Program) Pos() token.Pos           { return s.ValuePos }
func (s *LetStmt) Pos() token.Pos           { return s.ValuePos }
func (s *ExpressionStmt) Pos() token.Pos    { return s.ValuePos }
func (s *BlockStmt) Pos() token.Pos         { return s.ValuePos }
func (s *CommentStmt) Pos() token.Pos       { return s.ValuePos }
func (s *ConditionalStmt) Pos() token.Pos   { return s.ValuePos }
func (s *ForLoopStmt) Pos() token.Pos       { return s.ValuePos }
func (s *StringLiteral) Pos() token.Pos     { return s.ValuePos }
func (s *FunctionLiteral) Pos() token.Pos   { return s.ValuePos }
func (s *IntegerLiteral) Pos() token.Pos    { return s.ValuePos }
func (s *ArrayLiteral) Pos() token.Pos      { return s.ValuePos }
func (s *HashLiteral) Pos() token.Pos       { return s.ValuePos }
func (s *BooleanLiteral) Pos() token.Pos    { return s.ValuePos }
func (s *FloatLiteral) Pos() token.Pos      { return s.ValuePos }
func (s *Identifier) Pos() token.Pos        { return s.ValuePos }
func (s *ReassignmentStmt) Pos() token.Pos  { return s.ValuePos }
func (s *ErrorStmt) Pos() token.Pos         { return s.ValuePos }
func (s *InfixExpression) Pos() token.Pos   { return s.ValuePos }
func (s *IfStmt) Pos() token.Pos            { return s.ValuePos }
func (s *PrefixExpression) Pos() token.Pos  { return s.ValuePos }
func (s *ReturnExpression) Pos() token.Pos  { return s.ValuePos }
func (s *PostfixExpression) Pos() token.Pos { return s.ValuePos }
func (s *CallExpression) Pos() token.Pos    { return s.ValuePos }
func (s *IndexExpression) Pos() token.Pos   { return s.ValuePos }

func (s *Program) Literal() string          { return "" } // TODO
func (s *LetStmt) Literal() string          { return s.Token.Literal }
func (s *ExpressionStmt) Literal() string   { return s.Token.Literal }
func (s *BlockStmt) Literal() string        { return "" } // TODO
func (s *CommentStmt) Literal() string      { return "" } // TODO
func (s *ConditionalStmt) Literal() string  { return "" } // TODO
func (s *ForLoopStmt) Literal() string      { return s.Token.Literal }
func (s *StringLiteral) Literal() string    { return s.Value }
func (s *FunctionLiteral) Literal() string  { return s.Token.Literal } //TODO
func (s *IntegerLiteral) Literal() string   { return fmt.Sprint(s.Value) }
func (s *ArrayLiteral) Literal() string     { return "" } // TODO
func (s *HashLiteral) Literal() string      { return "" } // TODO
func (s *BooleanLiteral) Literal() string   { return fmt.Sprint(s.Value) }
func (s *FloatLiteral) Literal() string     { return fmt.Sprint(s.Value) }
func (s *Identifier) Literal() string       { return s.Value }
func (s *ReassignmentStmt) Literal() string { return s.Name.Literal() }
func (s *ErrorStmt) Literal() string        { return s.Value }
func (s *InfixExpression) Literal() string {
	return fmt.Sprintf("(%s %s %s)", s.Left.Literal(), s.Operator, s.Right.Literal())
}
func (s *IfStmt) Literal() string { return s.Token.Literal }
func (s *PrefixExpression) Literal() string {
	return fmt.Sprintf("%s%s", s.Token.Literal, s.Right.Literal())
}
func (s *ReturnExpression) Literal() string {
	return fmt.Sprintf("%s%s", s.Token.Literal, s.Expr.Literal())
}
func (s *PostfixExpression) Literal() string { return s.Token.Literal }
func (s *CallExpression) Literal() string    { return s.Token.Literal }
func (s *IndexExpression) Literal() string   { return s.Token.Literal }

func (s *Program) TokenType() token.TokenType           { return s.Token.Type }
func (s *LetStmt) TokenType() token.TokenType           { return s.Token.Type }
func (s *ExpressionStmt) TokenType() token.TokenType    { return s.Token.Type }
func (s *BlockStmt) TokenType() token.TokenType         { return s.Token.Type }
func (s *CommentStmt) TokenType() token.TokenType       { return s.Token.Type }
func (s *ForLoopStmt) TokenType() token.TokenType       { return s.Token.Type }
func (s *ConditionalStmt) TokenType() token.TokenType   { return s.Token.Type }
func (s *StringLiteral) TokenType() token.TokenType     { return s.Token.Type }
func (s *FunctionLiteral) TokenType() token.TokenType   { return s.Token.Type }
func (s *IntegerLiteral) TokenType() token.TokenType    { return s.Token.Type }
func (s *ArrayLiteral) TokenType() token.TokenType      { return s.Token.Type }
func (s *HashLiteral) TokenType() token.TokenType       { return s.Token.Type }
func (s *BooleanLiteral) TokenType() token.TokenType    { return s.Token.Type }
func (s *FloatLiteral) TokenType() token.TokenType      { return s.Token.Type }
func (s *Identifier) TokenType() token.TokenType        { return s.Token.Type }
func (s *ReassignmentStmt) TokenType() token.TokenType  { return s.Token.Type }
func (s *ErrorStmt) TokenType() token.TokenType         { return s.Token.Type }
func (s *IfStmt) TokenType() token.TokenType            { return s.Token.Type }
func (s *InfixExpression) TokenType() token.TokenType   { return s.Token.Type }
func (s *PrefixExpression) TokenType() token.TokenType  { return s.Token.Type }
func (s *ReturnExpression) TokenType() token.TokenType  { return s.Token.Type }
func (s *PostfixExpression) TokenType() token.TokenType { return s.Token.Type }
func (s *CallExpression) TokenType() token.TokenType    { return s.Token.Type }
func (s *IndexExpression) TokenType() token.TokenType   { return s.Token.Type }
