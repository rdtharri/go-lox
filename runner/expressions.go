package runner

import (
	"fmt"
	"strings"
)

type ExpressionVisitor interface {
	VisitBinaryExpression(*BinaryExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitUnaryExpression(*UnaryExpression)
}

type PrinterVistor struct {
	Messages []string
	Depth    int
}

func (pv *PrinterVistor) offset() string {
	return strings.Repeat(" ", pv.Depth*4)
}

func (pv *PrinterVistor) appendMessage(app string) {
	pv.Messages = append(pv.Messages,app)
}

func (pv *PrinterVistor) paren(leader string, args ...Expression) {
	pv.appendMessage("("+leader+" ")
	for _, arg := range args {
		arg.Accept(pv)
	}
	pv.appendMessage(" )")
}

func (pv *PrinterVistor) VisitBinaryExpression(be *BinaryExpression) {
	pv.paren(be.Operator.Lexeme, be.Left, be.Right)
}

func (pv *PrinterVistor) VisitGroupingExpression(ge *GroupingExpression) {
	pv.paren("", ge.Expression)
}

func (pv *PrinterVistor) VisitLiteralExpression(le *LiteralExpression) {
	pv.appendMessage(" " + le.Token.Lexeme + " ")
}

func (pv *PrinterVistor) VisitUnaryExpression(ue *UnaryExpression) {
	pv.paren(ue.Operator.Lexeme, ue.Right)
}

func (pv *PrinterVistor) print(exp Expression) {
	exp.Accept(pv)
	for _, message := range pv.Messages {
		fmt.Print(message)
	}
}

type Expression interface {
	Accept(ExpressionVisitor)
}

type BinaryExpression struct {
	Operator Token
	Left     Expression
	Right    Expression
}

func (be *BinaryExpression) Accept(v ExpressionVisitor) {
	v.VisitBinaryExpression(be)
}

type GroupingExpression struct {
	Expression Expression
}

func (ge *GroupingExpression) Accept(v ExpressionVisitor) {
	v.VisitGroupingExpression(ge)
}

type UnaryExpression struct {
	Operator Token
	Right    Expression
}

func (ge *UnaryExpression) Accept(v ExpressionVisitor) {
	v.VisitUnaryExpression(ge)
}

type LiteralExpression struct {
	Token Token
}

func (ge *LiteralExpression) Accept(v ExpressionVisitor) {
	v.VisitLiteralExpression(ge)
}
