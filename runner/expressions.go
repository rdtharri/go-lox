package runner

type ExpressionVisitor interface {
	VisitBinaryExpression(*BinaryExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitUnaryExpression(*UnaryExpression)
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
