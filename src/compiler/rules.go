package compiler

import (
	. "../token"
)

type ParseFn func(bool)
type ParseRule struct {
	prefix ParseFn
	infix  ParseFn
	postfix ParseFn
	Prec   Precedence
}

func (c *Compiler) LoadRules() {
	c.Rules = []ParseRule{
		{c.Grouping, c.Call,nil, PREC_NONE},    // TOKEN_LEFT_PAREN
		{nil, nil, nil,PREC_NONE},          // TOKEN_RIGHT_PAREN
		{nil, nil, nil,PREC_NONE},          // TOKEN_LEFT_BRACE
		{nil, nil, nil,PREC_NONE},          // TOKEN_RIGHT_BRACE
		{nil, c.Index, nil,PREC_INDEX},       // TOKEN_LEFT_BRACKET
		{nil, nil, nil,PREC_NONE},          // TOKEN_RIGHT_BRACKET
		{nil, nil, nil,PREC_NONE},          // TOKEN_COMMA
		{nil, c.Dot, nil,PREC_CALL},          // TOKEN_DOT
		{c.Unary, c.Binary, nil,PREC_TERM},     // TOKEN_MINUS
		{nil, c.Binary, nil,PREC_TERM},       // TOKEN_PLUS
		{nil, nil, nil,PREC_NONE},          // TOKEN_SEMICOLON
	    // 10
		{nil, c.Binary, nil,PREC_FACTOR},     // TOKEN_SLASH
		{nil, c.Binary, nil,PREC_FACTOR},     // TOKEN_STAR
		{nil, nil, nil,PREC_NONE},          // TOKEN_AT
		{nil, nil, nil,PREC_NONE},          // TOKEN_CR
		{nil, nil, nil,PREC_NONE},          // TOKEN_COLON
		{nil, nil, nil,PREC_NONE},          // TOKEN_PERCENT
		{nil, nil, nil,PREC_NONE},          // TOKEN_TILDE
		{nil, nil, nil,PREC_NONE},          // TOKEN_QUESTION
		{nil, nil, nil,PREC_NONE},          // TOKEN_HAT
		{nil, nil, nil,PREC_NONE},          // TOKEN_DOLLAR
		// 20
		{nil, nil, nil,PREC_NONE},          // TOKEN_BAR
		{nil, nil, nil,PREC_NONE},          // TOKEN_BACKTICK
		{c.Unary, nil, nil,PREC_NONE},        // TOKEN_BANG
		{nil, c.Binary, nil,PREC_EQUALITY},   // TOKEN_BANG_EQUAL
		{nil, nil, nil,PREC_NONE},          // TOKEN_EQUAL
		{nil, c.Binary, nil,PREC_EQUALITY},   // TOKEN_EQUAL_EQUAL
		{nil, c.Binary, nil,PREC_COMPARISON}, // TOKEN_GREATER
		{nil, c.Binary, nil,PREC_COMPARISON}, // TOKEN_GREATER_EQUAL
		{nil, c.Binary, nil,PREC_COMPARISON}, // TOKEN_LESS
		{nil, c.Binary, nil,PREC_COMPARISON}, // TOKEN_LESS_EQUAL
		// 30
		{c.Unary, nil, c.Postary,PREC_INCR},       // TOKEN_PLUS_PLUS
		{nil, nil, c.Postary,PREC_NONE},       // TOKEN_MINUS_MINUS
		{c.Array, nil, nil,PREC_ARRAY},       // TOKEN_LEFT_ARRAY
		{c.HMap, nil, nil,PREC_LIST},         // TOKEN_LEFT_LIST
		{c.Variable, nil, nil,PREC_NONE},     // TOKEN_IDENTIFIER
		{c.String, nil, nil,PREC_NONE},       // TOKEN_STRING
		{c.Integer, nil, nil,PREC_NONE},      // TOKEN_INTEGER
		{c.Float, nil, nil,PREC_NONE},        // TOKEN_DECIMAL
		{nil, nil, nil,PREC_NONE},          // TOKEN_HMAP
		{nil, c.and_, nil,PREC_AND},          // TOKEN_AND
		// 40
		{nil, nil, nil,PREC_NONE},          // TOKEN_CLASS
		{nil, nil, nil,PREC_NONE},          // TOKEN_ELSE
		{c.Literal, nil, nil,PREC_NONE},      // TOKEN_FALSE
		{nil, nil, nil,PREC_NONE},          // TOKEN_FOR
		{nil, nil, nil,PREC_NONE},          // TOKEN_FUNC
		{nil, nil, nil,PREC_NONE},          // TOKEN_IF
		{c.Literal, nil, nil,PREC_NONE},      // TOKEN_NIL
		{nil, c.or_, nil,PREC_OR},            // TOKEN_OR
		{nil, nil, nil,PREC_NONE},          // TOKEN_PRINT
		{nil, nil, nil,PREC_NONE},          // TOKEN_RETURN
		// 50
		{nil, nil, nil,PREC_NONE},          // TOKEN_SUPER
		{nil, nil, nil,PREC_NONE},          // TOKEN_THIS
		{c.Literal, nil, nil,PREC_NONE},      // TOKEN_TRUE
		{nil, nil, nil,PREC_NONE},          // TOKEN_VAR
		{nil, nil, nil,PREC_NONE},          // TOKEN_WHILE
		{nil, nil, nil,PREC_NONE},          // TOKEN_ERROR
		{nil, nil, nil,PREC_NONE},          // TOKEN_EOF
		{nil, nil, nil,PREC_NONE},          // TOKEN_INCLUDE
		{nil, nil, nil,PREC_NONE},          // TOKEN_INTEGER
		{nil, nil, nil,PREC_NONE},          // TOKEN_FLOAT
		// 60
		{nil, nil, nil,PREC_NONE},          // TOKEN_BOOL
		{nil, nil, nil,PREC_NONE},          // TOKEN_STRING
		{nil, nil, nil,PREC_NONE},          // TOKEN_ARRAY
		{c.SqlSelect, nil, nil,PREC_NONE},    // TOKEN_SELECT
		{nil, nil, nil,PREC_NONE},    	  // TOKEN_INSERT
		{nil, nil, nil,PREC_NONE},          // TOKEN_UPDATE
		{nil, nil, nil,PREC_NONE},          // TOKEN_DELETE
		{nil, nil, nil,PREC_NONE},          // TOKEN_FROM
		{nil, nil, nil,PREC_NONE},          // TOKEN_JOIN
		{nil, nil, nil,PREC_NONE},          // TOKEN_LEFT
	    // 70
		{nil, nil, nil,PREC_NONE},          // TOKEN_RIGHT
		{nil, nil, nil,PREC_NONE},          // TOKEN_CROSSJOIN
		{nil, nil, nil,PREC_NONE},          // TOKEN_WHERE
		{nil, nil, nil,PREC_NONE},          // TOKEN_ALL
		{nil, nil, nil,PREC_NONE},          // TOKEN_ORDER
		{nil, nil, nil,PREC_NONE},          // TOKEN_GROUP
		{nil, nil, nil,PREC_NONE},          // TOKEN_BY
		{nil, nil, nil,PREC_NONE},          // TOKEN_INTO
		{nil, nil, nil,PREC_NONE},          // TOKEN_VALUES
		{nil, nil, nil,PREC_NONE},          // TOKEN_AS
		// 80
		{nil, nil, nil,PREC_NONE},          // TOKEN_ON
		{c.Browse, nil, nil,PREC_NONE},       // TOKEN_BROWSE
		{nil, nil, nil,PREC_NONE},          // TOKEN_CREATE
		{nil, nil, nil,PREC_NONE},          // TOKEN_TABLE
		{nil, nil, nil,PREC_NONE},          // TOKEN_COLUMN
		{nil, nil, nil,PREC_NONE},          // TOKEN_ROW
		{nil, nil, nil,PREC_NONE},          // TOKEN_VIEW
		{nil, nil, nil,PREC_NONE},          // TOKEN_HAVING
		{nil, nil, nil,PREC_NONE},          // TOKEN_DISTINCT
		{nil, nil, nil,PREC_NONE},          // TOKEN_TOP
		// 90
		{nil, nil, nil,PREC_NONE},          // TOKEN_STEP
		{nil, nil, nil,PREC_NONE},          // TOKEN_TO
	}
}
func (c *Compiler) GetRule(t_type TokenType) *ParseRule {
	return &c.Rules[t_type]
}
