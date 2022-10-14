package runner

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement)
	VisitPrintStatement(*PrintStatement)
	VisitVarStatement(*VarStatement)
	VisitBlockStatement(*BlockStatement)
	VisitIfStatement(*IfStatement)
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
	Name        Token
	Initializer Expression
}

func (vs *VarStatement) Accept(v StatementVisitor) {
	v.VisitVarStatement(vs)
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) Accept(v StatementVisitor) {
	v.VisitBlockStatement(bs)
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (is *IfStatement) Accept(v StatementVisitor) {
	v.VisitIfStatement(is)
}
