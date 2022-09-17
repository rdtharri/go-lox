package runner

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement)
	VisitPrintStatement(*PrintStatement)
	VisitVarStatement(*VarStatement)
}

type Statement interface {
	Accept(StatementVisitor)
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) Accept(v StatementVisitor) {
	v.VisitExpressionStatement(es)
}

type PrintStatement struct {
	Expression Expression
}

func (ps *PrintStatement) Accept(v StatementVisitor) {
	v.VisitPrintStatement(ps)
}

type VarStatement struct {
	Name Token
	Initializer Expression
}

func (vs *VarStatement) Accept(v StatementVisitor) {
	v.VisitVarStatement(vs)
}
