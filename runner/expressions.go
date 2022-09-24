package runner

type ExpressionVisitor interface {
	VisitBinaryExpression(*BinaryExpression) interface{}
	VisitGroupingExpression(*GroupingExpression) interface{}
	VisitLiteralExpression(*LiteralExpression) interface{}
	VisitUnaryExpression(*UnaryExpression) interface{}
	VisitVarExpression(*VarExpression) interface{}
	VisitAssignExpression(*AssignExpression) interface{}
}

type Expression interface {
	Accept(ExpressionVisitor) interface{}
}

type BinaryExpression struct {
	Operator Token
	Left     Expression
	Right    Expression
}

func (be *BinaryExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitBinaryExpression(be)
}

type GroupingExpression struct {
	Expression Expression
}

func (ge *GroupingExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitGroupingExpression(ge)
}

type UnaryExpression struct {
	Operator Token
	Right    Expression
}

func (ge *UnaryExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitUnaryExpression(ge)
}

type LiteralExpression struct {
	Token Token
}

func (ge *LiteralExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitLiteralExpression(ge)
}

type VarExpression struct {
	Name Token
}

func (ve *VarExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitVarExpression(ve)
}

type AssignExpression struct {
	Name Token
	Value Expression
}

func (ae *AssignExpression) Accept(v ExpressionVisitor) interface{} {
	return v.VisitAssignExpression(ae)
}
