package main

type ParseFn func(bool)
type ParseRule struct {
	prefix  ParseFn
	infix   ParseFn
	postfix ParseFn
	Prec    Precedence
}

func (c *Compiler) LoadRules() {
	c.Rules = []ParseRule{
		{c.Grouping, c.Call, nil, PREC_CALL}, // TOKEN_LEFT_PAREN
		{nil, nil, nil, PREC_NONE},           // TOKEN_RIGHT_PAREN
		{nil, nil, nil, PREC_NONE},           // TOKEN_LEFT_BRACE
		{nil, nil, nil, PREC_NONE},           // TOKEN_RIGHT_BRACE
		{c.Index, nil, nil, PREC_NONE},       // TOKEN_LEFT_BRACKET
		{nil, nil, nil, PREC_NONE},           // TOKEN_RIGHT_BRACKET
		{nil, nil, nil, PREC_NONE},           // TOKEN_COMMA
		{nil, nil, nil, PREC_CALL},           // TOKEN_DOT
		{c.Unary, c.Binary, nil, PREC_TERM},  // TOKEN_MINUS
		{nil, c.Binary, nil, PREC_TERM},      // TOKEN_PLUS
		{nil, nil, nil, PREC_NONE},           // TOKEN_SEMICOLON
		// 10
		{nil, c.Binary, nil, PREC_FACTOR}, // TOKEN_SLASH
		{nil, c.Binary, nil, PREC_FACTOR}, // TOKEN_STAR
		{nil, nil, nil, PREC_NONE},        // TOKEN_AT
		{nil, nil, nil, PREC_NONE},        // TOKEN_CR
		{nil, nil, nil, PREC_NONE},        // TOKEN_COLON
		{nil, nil, nil, PREC_NONE},        // TOKEN_PERCENT
		{nil, nil, nil, PREC_NONE},        // TOKEN_TILDE
		{nil, nil, nil, PREC_NONE},        // TOKEN_QUESTION
		{nil, c.Binary, nil, PREC_FACTOR}, // TOKEN_HAT
		{c.Dollar, nil, nil, PREC_NONE},        // TOKEN_DOLLAR
		// 20
		{nil, nil, nil, PREC_NONE},            // TOKEN_BAR
		{nil, nil, nil, PREC_NONE},            // TOKEN_BACKTICK
		{c.Unary, nil, nil, PREC_NONE},        // TOKEN_BANG
		{nil, c.Binary, nil, PREC_EQUALITY},   // TOKEN_BANG_EQUAL
		{nil, nil, nil, PREC_NONE},            // TOKEN_EQUAL
		{nil, c.Binary, nil, PREC_EQUALITY},   // TOKEN_EQUAL_EQUAL
		{nil, c.Binary, nil, PREC_COMPARISON}, // TOKEN_GREATER
		{nil, c.Binary, nil, PREC_COMPARISON}, // TOKEN_GREATER_EQUAL
		{nil, c.Binary, nil, PREC_COMPARISON}, // TOKEN_LESS
		{nil, c.Binary, nil, PREC_COMPARISON}, // TOKEN_LESS_EQUAL
		// 30
		{c.Unary, nil, c.Postary, PREC_INCR}, // TOKEN_PLUS_PLUS
		{nil, nil, c.Postary, PREC_NONE},     // TOKEN_MINUS_MINUS
		{c.Array, nil, nil, PREC_NONE},      // TOKEN_ARRAY
		{c.List, nil, nil, PREC_LIST},        // TOKEN_LIST
		{c.Variable, nil, nil, PREC_NONE},    // TOKEN_IDENTIFIER
		{c.String, nil, nil, PREC_NONE},      // TOKEN_STRING
		{c.Integer, nil, nil, PREC_NONE},     // TOKEN_INTEGER
		{c.Float, nil, nil, PREC_NONE},       // TOKEN_DECIMAL
		{nil, nil, nil, PREC_NONE},           // TOKEN_HMAP
		{nil, c.and_, nil, PREC_AND},         // TOKEN_AND
		// 40
		{c.Class, nil, nil, PREC_NONE},    // TOKEN_CLASS
		{nil, nil, nil, PREC_NONE},        // TOKEN_ELSE
		{c.Boolean, nil, nil, PREC_NONE},  // TOKEN_FALSE
		{nil, nil, nil, PREC_NONE},        // TOKEN_FOR
		{c.Function, nil, nil, PREC_NONE}, // TOKEN_FUNC
		{nil, nil, nil, PREC_NONE},        // TOKEN_IF
		{c.Literal, nil, nil, PREC_NONE},  // TOKEN_NIL
		{nil, c.or_, nil, PREC_OR},        // TOKEN_OR
		{nil, nil, nil, PREC_NONE},        // TOKEN_RETURN
		{nil, nil, nil, PREC_NONE},        // TOKEN_SUPER
		// 50
		{c.Enum, nil, nil, PREC_NONE},    // TOKEN_ENUM
		{c.Boolean, nil, nil, PREC_NONE}, // TOKEN_TRUE
		{nil, nil, nil, PREC_NONE},       // TOKEN_VAR
		{nil, nil, nil, PREC_NONE},       // TOKEN_WHILE
		{nil, nil, nil, PREC_NONE},       // TOKEN_ERROR
		{nil, nil, nil, PREC_NONE},       // TOKEN_EOF
		{nil, nil, nil, PREC_NONE},       // TOKEN_INCLUDE
		{c.IntegerType, nil, nil, PREC_NONE},       // TOKEN_TYPE_INTEGER
		{c.FloatType, nil, nil, PREC_NONE},       // TOKEN_TYPE_FLOAT
		{c.BoolType, nil, nil, PREC_NONE},       // TOKEN_TYPE_BOOL
		// 60
		{c.StringType, nil, nil, PREC_NONE},         // TOKEN_TYPE_STRING
		{c._array, nil, nil, PREC_NONE},    // TOKEN_TYPE_ARRAY
		{c.SqlSelect, nil, nil, PREC_NONE}, // TOKEN_SELECT
		{nil, nil, nil, PREC_NONE},         // TOKEN_INSERT
		{nil, nil, nil, PREC_NONE},         // TOKEN_UPDATE
		{nil, nil, nil, PREC_NONE},         // TOKEN_DELETE
		{nil, nil, nil, PREC_NONE},         // TOKEN_FROM
		{nil, nil, nil, PREC_NONE},         // TOKEN_JOIN
		{nil, nil, nil, PREC_NONE},         // TOKEN_LEFT
		{nil, nil, nil, PREC_NONE},         // TOKEN_RIGHT
		// 70
		{nil, nil, nil, PREC_NONE}, // TOKEN_CROSSJOIN
		{nil, nil, nil, PREC_NONE}, // TOKEN_WHERE
		{nil, nil, nil, PREC_NONE}, // TOKEN_ALL
		{nil, nil, nil, PREC_NONE}, // TOKEN_ORDER
		{nil, nil, nil, PREC_NONE}, // TOKEN_GROUP
		{nil, nil, nil, PREC_NONE}, // TOKEN_BY
		{nil, nil, nil, PREC_NONE}, // TOKEN_INTO
		{nil, nil, nil, PREC_NONE}, // TOKEN_VALUES
		{nil, nil, nil, PREC_NONE}, // TOKEN_AS
		{nil, nil, nil, PREC_NONE}, // TOKEN_ON
		// 80
		{c.Browse, nil, nil, PREC_NONE}, // TOKEN_BROWSE
		{nil, nil, nil, PREC_NONE},      // TOKEN_CREATE
		{nil, nil, nil, PREC_NONE},      // TOKEN_TABLE
		{nil, nil, nil, PREC_NONE},      // TOKEN_COLUMN
		{nil, nil, nil, PREC_NONE},      // TOKEN_ROW
		{nil, nil, nil, PREC_NONE},      // TOKEN_VIEW
		{nil, nil, nil, PREC_NONE},      // TOKEN_HAVING
		{nil, nil, nil, PREC_NONE},      // TOKEN_DISTINCT
		{nil, nil, nil, PREC_NONE},      // TOKEN_TOP
		{nil, nil, nil, PREC_NONE},      // TOKEN_STEP
		// 90
		{nil, nil, nil, PREC_NONE}, // TOKEN_TO
		{nil, nil, nil, PREC_NONE}, // TOKEN_WHEN
		{nil, nil, nil, PREC_NONE}, // TOKEN_CASE
		{nil, nil, nil, PREC_NONE}, // TOKEN_DEFAULT
		{nil, nil, nil, PREC_NONE}, // TOKEN_SWITCH
		{nil, nil, nil, PREC_NONE}, // TOKEN_BREAK
		{nil, nil, nil, PREC_NONE}, // TOKEN_CONTINUE
		{nil, nil, nil, PREC_NONE}, // TOKEN_SCAN
		{nil, nil, nil, PREC_NONE}, // TOKEN_PROPERTY
		{nil, nil, nil, PREC_NONE}, // TOKEN_METHOD
		//100
		{nil, nil, nil, PREC_NONE}, // TOKEN_PRIVATE
		{nil, nil, nil, PREC_NONE}, // TOKEN_PUBLIC
		{nil, nil, nil, PREC_NONE}, // TOKEN_PROTECTED
		{c.NewList, nil, nil, PREC_NONE}, // TOKEN_LIST_TYPE
		{nil, nil, nil, PREC_NONE}, // TOKEN_INDEX
		{c.New, nil, nil, PREC_NONE}, // TOKEN_NEW
		{nil, nil, nil, PREC_NONE}, // TOKEN_NOT
		{nil, nil, nil, PREC_NONE}, // TOKEN_NULL
		{nil, nil, nil, PREC_NONE}, // TOKEN_STRING2
		{nil, nil, nil, PREC_NONE}, // TOKEN_BEGIN_VAR
		{nil, nil, nil, PREC_NONE}, // TOKEN_END_VAR

		// SQL commands
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ABORT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ACTION
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_AFTER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ALL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ALTER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ALWAYS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ANALYZE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_AND
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_AS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ASC

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ATTACH
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_AUTOINCREMENT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_BEFORE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_BEGIN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_BETWEEN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_BY
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CASCADE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CASE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CAST
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CHECK

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_COLLATE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_COLUMN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_COMMIT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CONFLICT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CONSTRAINT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CREATE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CROSS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CURRENT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CURRENT_DATE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CURRENT_TIME

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_CURRENT_TIMESTAMP
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DATABASE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DEFAULT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DEFERRABLE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DEFERRED
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DELETE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DESC
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DETACH
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DISTINCT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DO
		//
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_DROP
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EACH
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ELSE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_END
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ESCAPE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EXCEPT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EXCLUDE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EXCLUSIVE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EXISTS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_EXPLAIN
		//
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FAIL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FILTER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FIRST
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FOLLOWING
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FOR
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FOREIGN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FROM
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_FULL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_GENERATED
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_GLOB
		//
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_GROUP
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_GROUPS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_HAVING
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_IF
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_IGNORE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_IMMEDIATE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_IN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INDEX
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INDEXED
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INITIALLY

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INNER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INSERT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INSTEAD
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INTERSECT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_INTO
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_IS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ISNULL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_JOIN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_KEY
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_LAST

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_LEFT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_LIKE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_LIMIT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_MATCH
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NATURAL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NO
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NOT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NOTHING
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NOTNULL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NULL

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_NULLS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OF
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OFFSET
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ON
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OR
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ORDER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OTHERS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OUTER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_OVER
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_PARTITION

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_PLAN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_PRAGMA
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_PRECEDING
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_PRIMARY
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_QUERY
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RAISE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RANGE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RECURSIVE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_REFERENCES
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_REGEXP

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_REINDEX
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RELEASE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RENAME
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_REPLACE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RESTRICT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_RIGHT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ROLLBACK
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ROW
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_ROWS
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_SAVEPOINT

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_SELECT
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_SET
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TABLE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TEMP
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TEMPORARY
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_THEN
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TIES
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TO
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TRANSACTION
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_TRIGGER

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_UNBOUNDED
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_UNION
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_UNIQUE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_UPDATE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_USING
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_VACUUM
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_VALUES
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_VIEW
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_VIRTUAL
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_WHEN

		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_WHERE
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_WINDOW
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_WITH
		{nil, nil, nil, PREC_NONE}, //TOKEN_SQL_WITHOUT
		{nil, nil, nil, PREC_NONE}, //TOKEN_DOUBLE_COLON


	}
}
func (c *Compiler) GetRule(t_type TokenType) *ParseRule {
	return &c.Rules[t_type]
}
