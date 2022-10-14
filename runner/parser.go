package runner

import (
	"errors"
	"fmt"
)

type Parser struct {
	Tokens  []Token
	current int
	Runner  *LoxRunner
}

func NewParser(tokens []Token, runner *LoxRunner) *Parser {
	parser := new(Parser)
	parser.Tokens = tokens
	parser.Runner = runner
	return parser
}

func (p *Parser) expression() Expression {
	return p.assignment()
}

func (p *Parser) assignment() Expression {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if varExp, ok := expr.(*VarExpression); ok {
			return &AssignExpression{
				Name:  varExp.Name,
				Value: value,
			}
		}

		p.error(equals, "Invalid assignment target")
	}

	return expr
}

func (p *Parser) or() Expression {
	expr := p.and()

	for p.match(OR) {
		expr = &LogicalExpression{
			Left:     expr,
			Right:    p.and(),
			Operator: p.previous(),
		}
	}

	return expr
}

func (p *Parser) and() Expression {
	expr := p.equality()

	for p.match(AND) {
		expr = &LogicalExpression{
			Left:     expr,
			Right:    p.equality(),
			Operator: p.previous(),
		}
	}

	return expr
}

func (p *Parser) equality() Expression {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left:     expr,
			Right:    p.comparison(),
		}
	}

	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left:     expr,
			Right:    p.term(),
		}
	}

	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left:     expr,
			Right:    p.factor(),
		}
	}

	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left:     expr,
			Right:    p.unary(),
		}
	}

	return expr
}

func (p *Parser) unary() Expression {
	if p.match(BANG, MINUS) {
		return &UnaryExpression{
			Operator: p.previous(),
			Right:    p.unary(),
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expression {
	if p.match(FALSE, TRUE, NIL, NUMBER, STRING) {
		return &LiteralExpression{
			Token: p.previous(),
		}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpression{
			Expression: expr,
		}
	}

	if p.match(IDENTIFIER) {
		return &VarExpression{
			Name: p.previous(),
		}
	}

	panic(p.error(p.peek(), "Unexpected token"))
}

func (p *Parser) consume(ttype TokenType, message string) Token {
	if p.check(ttype) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) error(token Token, message string) error {
	p.Runner.tokenError(token, message)
	return errors.New(message)
}

func (p *Parser) match(checks ...TokenType) bool {
	for _, ttype := range checks {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == ttype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) parse() []Statement {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	statements := make([]Statement, 0)
	for !p.isAtEnd() {
		statements = append(
			statements,
			p.declaration(),
		)
	}
	return statements
}

func (p *Parser) declaration() Statement {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Statement {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expression
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after declaration")

	return &VarStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) block() []Statement {
	statements := make([]Statement, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(
			statements,
			p.declaration(),
		)
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (p *Parser) statement() Statement {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return &BlockStatement{
			Statements: p.block(),
		}
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() Statement {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Statement
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Statement {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &PrintStatement{
		Expression: value,
	}
}

func (p *Parser) expressionStatement() Statement {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &ExpressionStatement{
		Expression: value,
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
