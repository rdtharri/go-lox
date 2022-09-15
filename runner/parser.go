package runner

import (
	"errors"
)

type Parser struct {
	Tokens []Token
	current int
	Runner *LoxRunner
}

func NewParser(tokens []Token, runner *LoxRunner) *Parser {
	parser := new(Parser)
	parser.Tokens = tokens
	parser.Runner = runner
	return parser
}


func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left: expr,
			Right: p.comparison(),
		}
	}

	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left: expr,
			Right: p.term(),
		}
	}

	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left: expr,
			Right: p.factor(),
		}
	}

	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		expr = &BinaryExpression{
			Operator: p.previous(),
			Left: expr,
			Right: p.unary(),
		}
	}

	return expr
}

func (p *Parser) unary() Expression {
	if p.match(BANG, MINUS) {
		return &UnaryExpression{
			Operator: p.previous(),
			Right: p.unary(),
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


	panic(p.error(p.peek(),"Unexpected token"))
}

func (p *Parser) consume(ttype TokenType, message string) Token {
	if p.check(ttype) {
		return p.advance()
	}

	panic(p.error(p.peek(),message))
}

func (p *Parser) error(token Token, message string) error {
	p.Runner.tokenError(token, message)
	return  errors.New(message)
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
	return p.Tokens[p.current - 1]
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

func (p *Parser) parse() Expression {
	return p.expression()
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
