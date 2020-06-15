package main

type TokenType int

type Token struct {
	Value  []byte
	Type   TokenType
	Start  int
	Length int
	Line   int
}

func (t *Token) ToString() string {
	return string(t.Value)
}

const (
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	// 10
	TOKEN_SLASH
	TOKEN_STAR
	TOKEN_AT
	TOKEN_CR
	TOKEN_COLON
	TOKEN_PERCENT
	TOKEN_TILDE
	TOKEN_QUESTION
	TOKEN_HAT
	TOKEN_DOLLAR
	// 20
	TOKEN_BAR
	TOKEN_BACKTICK
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	// 30
	TOKEN_PLUS_PLUS
	TOKEN_MINUS_MINUS
	TOKEN_ARRAY
	TOKEN_LIST
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_INTEGER
	TOKEN_DECIMAL
	TOKEN_HMAP
	TOKEN_AND
	// 40
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUNC
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_RETURN
	TOKEN_SUPER
	// 50
	TOKEN_ENUM
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE
	TOKEN_ERROR
	TOKEN_EOF
	TOKEN_INCLUDE
	TOKEN_TYPE_INTEGER
	TOKEN_TYPE_FLOAT
	TOKEN_TYPE_BOOL
	// 60
	TOKEN_TYPE_STRING
	TOKEN_TYPE_ARRAY
	// SQL Specific commands
	TOKEN_SELECT
	TOKEN_INSERT
	TOKEN_UPDATE
	TOKEN_DELETE
	TOKEN_FROM
	TOKEN_JOIN
	TOKEN_LEFT
	TOKEN_RIGHT
	// 70
	TOKEN_CROSSJOIN
	TOKEN_WHERE
	TOKEN_ALL
	TOKEN_ORDER
	TOKEN_GROUP
	TOKEN_BY
	TOKEN_INTO
	TOKEN_VALUES
	TOKEN_AS
	TOKEN_ON
	// 80
	TOKEN_BROWSE
	TOKEN_CREATE
	TOKEN_TABLE
	TOKEN_COLUMN
	TOKEN_ROW
	TOKEN_VIEW
	TOKEN_HAVING
	TOKEN_DISTINCT
	TOKEN_TOP
	TOKEN_STEP
	//90
	TOKEN_TO
	TOKEN_WHEN
	TOKEN_CASE
	TOKEN_DEFAULT
	TOKEN_SWITCH
	TOKEN_BREAK
	TOKEN_CONTINUE
	TOKEN_SCAN
	TOKEN_PROPERTY
	TOKEN_METHOD
	// 100
	TOKEN_PRIVATE
	TOKEN_PUBLIC
	TOKEN_PROTECTED
	TOKEN_LIST_TYPE
	TOKEN_INDEX
	TOKEN_NEW
	TOKEN_NOT
	TOKEN_NULL
	TOKEN_STRING2
	TOKEN_BEGIN_VAR
	TOKEN_END_VAR
	// SQL Specific tokens
	//
	TOKEN_SQL_ABORT
	TOKEN_SQL_ACTION
	TOKEN_SQL_AFTER
	TOKEN_SQL_ALL
	TOKEN_SQL_ALTER
	TOKEN_SQL_ALWAYS
	TOKEN_SQL_ANALYZE
	TOKEN_SQL_AND
	TOKEN_SQL_AS
	TOKEN_SQL_ASC
	//
	TOKEN_SQL_ATTACH
	TOKEN_SQL_AUTOINCREMENT
	TOKEN_SQL_BEFORE
	TOKEN_SQL_BEGIN
	TOKEN_SQL_BETWEEN
	TOKEN_SQL_BY
	TOKEN_SQL_CASCADE
	TOKEN_SQL_CASE
	TOKEN_SQL_CAST
	TOKEN_SQL_CHECK
	//
	TOKEN_SQL_COLLATE
	TOKEN_SQL_COLUMN
	TOKEN_SQL_COMMIT
	TOKEN_SQL_CONFLICT
	TOKEN_SQL_CONSTRAINT
	TOKEN_SQL_CREATE
	TOKEN_SQL_CROSS
	TOKEN_SQL_CURRENT
	TOKEN_SQL_CURRENT_DATE
	TOKEN_SQL_CURRENT_TIME
	//
	TOKEN_SQL_CURRENT_TIMESTAMP
	TOKEN_SQL_DATABASE
	TOKEN_SQL_DEFAULT
	TOKEN_SQL_DEFERRABLE
	TOKEN_SQL_DEFERRED
	TOKEN_SQL_DELETE
	TOKEN_SQL_DESC
	TOKEN_SQL_DETACH
	TOKEN_SQL_DISTINCT
	TOKEN_SQL_DO
	//
	TOKEN_SQL_DROP
	TOKEN_SQL_EACH
	TOKEN_SQL_ELSE
	TOKEN_SQL_END
	TOKEN_SQL_ESCAPE
	TOKEN_SQL_EXCEPT
	TOKEN_SQL_EXCLUDE
	TOKEN_SQL_EXCLUSIVE
	TOKEN_SQL_EXISTS
	TOKEN_SQL_EXPLAIN
	//
	TOKEN_SQL_FAIL
	TOKEN_SQL_FILTER
	TOKEN_SQL_FIRST
	TOKEN_SQL_FOLLOWING
	TOKEN_SQL_FOR
	TOKEN_SQL_FOREIGN
	TOKEN_SQL_FROM
	TOKEN_SQL_FULL
	TOKEN_SQL_GENERATED
	TOKEN_SQL_GLOB
	//
	TOKEN_SQL_GROUP
	TOKEN_SQL_GROUPS
	TOKEN_SQL_HAVING
	TOKEN_SQL_IF
	TOKEN_SQL_IGNORE
	TOKEN_SQL_IMMEDIATE
	TOKEN_SQL_IN
	TOKEN_SQL_INDEX
	TOKEN_SQL_INDEXED
	TOKEN_SQL_INITIALLY
	//
	TOKEN_SQL_INNER
	TOKEN_SQL_INSERT
	TOKEN_SQL_INSTEAD
	TOKEN_SQL_INTERSECT
	TOKEN_SQL_INTO
	TOKEN_SQL_IS
	TOKEN_SQL_ISNULL
	TOKEN_SQL_JOIN
	TOKEN_SQL_KEY
	TOKEN_SQL_LAST
	//
	TOKEN_SQL_LEFT
	TOKEN_SQL_LIKE
	TOKEN_SQL_LIMIT
	TOKEN_SQL_MATCH
	TOKEN_SQL_NATURAL
	TOKEN_SQL_NO
	TOKEN_SQL_NOT
	TOKEN_SQL_NOTHING
	TOKEN_SQL_NOTNULL
	TOKEN_SQL_NULL
	//
	TOKEN_SQL_NULLS
	TOKEN_SQL_OF
	TOKEN_SQL_OFFSET
	TOKEN_SQL_ON
	TOKEN_SQL_OR
	TOKEN_SQL_ORDER
	TOKEN_SQL_OTHERS
	TOKEN_SQL_OUTER
	TOKEN_SQL_OVER
	TOKEN_SQL_PARTITION
	//
	TOKEN_SQL_PLAN
	TOKEN_SQL_PRAGMA
	TOKEN_SQL_PRECEDING
	TOKEN_SQL_PRIMARY
	TOKEN_SQL_QUERY
	TOKEN_SQL_RAISE
	TOKEN_SQL_RANGE
	TOKEN_SQL_RECURSIVE
	TOKEN_SQL_REFERENCES
	TOKEN_SQL_REGEXP
	//
	TOKEN_SQL_REINDEX
	TOKEN_SQL_RELEASE
	TOKEN_SQL_RENAME
	TOKEN_SQL_REPLACE
	TOKEN_SQL_RESTRICT
	TOKEN_SQL_RIGHT
	TOKEN_SQL_ROLLBACK
	TOKEN_SQL_ROW
	TOKEN_SQL_ROWS
	TOKEN_SQL_SAVEPOINT
	//
	TOKEN_SQL_SELECT
	TOKEN_SQL_SET
	TOKEN_SQL_TABLE
	TOKEN_SQL_TEMP
	TOKEN_SQL_TEMPORARY
	TOKEN_SQL_THEN
	TOKEN_SQL_TIES
	TOKEN_SQL_TO
	TOKEN_SQL_TRANSACTION
	TOKEN_SQL_TRIGGER
	//
	TOKEN_SQL_UNBOUNDED
	TOKEN_SQL_UNION
	TOKEN_SQL_UNIQUE
	TOKEN_SQL_UPDATE
	TOKEN_SQL_USING
	TOKEN_SQL_VACUUM
	TOKEN_SQL_VALUES
	TOKEN_SQL_VIEW
	TOKEN_SQL_VIRTUAL
	TOKEN_SQL_WHEN
	//
	TOKEN_SQL_WHERE
	TOKEN_SQL_WINDOW
	TOKEN_SQL_WITH
	TOKEN_SQL_WITHOUT
	TOKEN_MODULE
	TOKEN_IMPORT
	TOKEN_DOUBLE_COLON
)

type TokenProperties struct {
	Type            TokenType
	IsCaseSensitive bool
}

var TokenLabels = map[string]TokenProperties{
	"(": {TOKEN_LEFT_PAREN, true},
	")": {TOKEN_RIGHT_PAREN, true},
	"{": {TOKEN_LEFT_BRACE, true},
	"}": {TOKEN_RIGHT_BRACE, true},
	"[": {TOKEN_LEFT_BRACKET, true},
	"]": {TOKEN_RIGHT_BRACKET, true},
	",": {TOKEN_COMMA, true},
	".": {TOKEN_DOT, true},
	"-": {TOKEN_MINUS, true},
	"+": {TOKEN_PLUS, true},
	";": {TOKEN_SEMICOLON, true},
	// 10
	"/":  {TOKEN_SLASH, true},
	"*":  {TOKEN_STAR, true},
	"@":  {TOKEN_AT, true},
	"\n": {TOKEN_CR, true},
	":":  {TOKEN_COLON, true},
	"%":  {TOKEN_PERCENT, true},
	"~":  {TOKEN_TILDE, true},
	"?":  {TOKEN_QUESTION, true},
	"^":  {TOKEN_HAT, true},
	"$":  {TOKEN_DOLLAR, true},
	// 20
	"|": {TOKEN_BAR, true},
	"`": {TOKEN_BACKTICK, true},
	// One or two character tokens.
	"!":  {TOKEN_BANG, true},
	"!=": {TOKEN_BANG_EQUAL, true},
	"=":  {TOKEN_EQUAL, true},
	"==": {TOKEN_EQUAL_EQUAL, true},
	">":  {TOKEN_GREATER, true},
	">=": {TOKEN_GREATER_EQUAL, true},
	"<":  {TOKEN_LESS, true},
	"<=": {TOKEN_LESS_EQUAL, true},
	// 30
	"++":     {TOKEN_PLUS_PLUS, true},
	"--":     {TOKEN_MINUS_MINUS, true},
	"@[":     {TOKEN_ARRAY, true},
	"@{":     {TOKEN_LIST, true},
	"and":    {TOKEN_AND, true},
	"class":  {TOKEN_CLASS, true},
	"else":   {TOKEN_ELSE, true},
	"false":  {TOKEN_FALSE, true},
	"for":    {TOKEN_FOR, true},
	"func":   {TOKEN_FUNC, true},
	"select": {TOKEN_SELECT, false},
	//40
	"if":     {TOKEN_IF, true},
	"nil":    {TOKEN_NIL, true},
	"or":     {TOKEN_OR, true},
	"return": {TOKEN_RETURN, true},
	"super":  {TOKEN_SUPER, true},
	"enum":   {TOKEN_ENUM, true},
	"true":   {TOKEN_TRUE, true},
	"var":    {TOKEN_VAR, true},
	"while":  {TOKEN_WHILE, true},
	"error":  {TOKEN_ERROR, true},
	// 50
	"include": {TOKEN_INCLUDE, true},
	"int":     {TOKEN_TYPE_INTEGER, true},
	"float":   {TOKEN_TYPE_FLOAT, true},
	"bool":    {TOKEN_TYPE_BOOL, true},
	"string":  {TOKEN_TYPE_STRING, true},
	"_array":  {TOKEN_TYPE_ARRAY, true},
	"insert":  {TOKEN_INSERT, false},
	"update":  {TOKEN_UPDATE, false},
	"delete":  {TOKEN_DELETE, false},
	// 60
	"from":      {TOKEN_FROM, true},
	"join":      {TOKEN_JOIN, true},
	"left":      {TOKEN_LEFT, true},
	"right":     {TOKEN_RIGHT, true},
	"crossjoin": {TOKEN_CROSSJOIN, true},
	"where":     {TOKEN_WHERE, true},
	"all":       {TOKEN_ALL, true},
	"order":     {TOKEN_ORDER, true},
	"group":     {TOKEN_GROUP, true},
	"by":        {TOKEN_BY, true},
	// 70
	"into":   {TOKEN_INTO, true},
	"values": {TOKEN_VALUES, true},
	"as":     {TOKEN_AS, true},
	"on":     {TOKEN_ON, true},
	"browse": {TOKEN_BROWSE, true},
	// DDL statements
	"create":   {TOKEN_CREATE, true},
	"table":    {TOKEN_TABLE, true},
	"column":   {TOKEN_COLUMN, true},
	"row":      {TOKEN_ROW, true},
	"view":     {TOKEN_VIEW, true},
	"having":   {TOKEN_HAVING, true},
	"distinct": {TOKEN_DISTINCT, true},
	"top":      {TOKEN_TOP, true},
	//
	"step":      {TOKEN_STEP, false},
	"to":        {TOKEN_TO, false},
	"when":      {TOKEN_WHEN, false},
	"case":      {TOKEN_CASE, false},
	"default":   {TOKEN_DEFAULT, false},
	"switch":    {TOKEN_SWITCH, false},
	"break":     {TOKEN_BREAK, false},
	"continue":  {TOKEN_CONTINUE, false},
	"scan":      {TOKEN_SCAN, false},
	"property":  {TOKEN_PROPERTY, false},
	"method":    {TOKEN_METHOD, false},
	"private":   {TOKEN_PRIVATE, false},
	"public":    {TOKEN_PUBLIC, false},
	"protected": {TOKEN_PROTECTED, false},
	"list":      {TOKEN_LIST_TYPE, true},
	"index":     {TOKEN_INDEX, false},
	"new":       {TOKEN_NEW, true},
	"not":       {TOKEN_NOT, false},
	"null":      {TOKEN_NULL, false},
	"<%":      {TOKEN_BEGIN_VAR, false},
	"%>":      {TOKEN_END_VAR, false},
	"module":      {TOKEN_MODULE, true},
	"import":      {TOKEN_IMPORT, true},
	"::":		   {TOKEN_DOUBLE_COLON, true},
}
var SqlTokenLabels = map[string]TokenProperties{
	// SQL Commnads
	"abort":{TOKEN_SQL_ABORT ,false},
	"action":{TOKEN_SQL_ACTION ,false},
	"after":{TOKEN_SQL_AFTER ,false},
	"all":{TOKEN_SQL_ALL ,false},
	"alter":{TOKEN_SQL_ALTER ,false},
	"always":{TOKEN_SQL_ALWAYS ,false},
	"analyze":{TOKEN_SQL_ANALYZE ,false},
	"and":{TOKEN_SQL_AND ,false},
	"as":{TOKEN_SQL_AS ,false},
	"asc":{TOKEN_SQL_ASC ,false},
	"attach":{TOKEN_SQL_ATTACH ,false},
	"autoincrement":{TOKEN_SQL_AUTOINCREMENT ,false},
	"before":{TOKEN_SQL_BEFORE ,false},
	"begin":{TOKEN_SQL_BEGIN ,false},
	"between":{TOKEN_SQL_BETWEEN ,false},
	"by":{TOKEN_SQL_BY ,false},
	"cascade":{TOKEN_SQL_CASCADE ,false},
	"case":{TOKEN_SQL_CASE ,false},
	"cast":{TOKEN_SQL_CAST ,false},
	"check":{TOKEN_SQL_CHECK ,false},
	"collate":{TOKEN_SQL_COLLATE ,false},
	"column":{TOKEN_SQL_COLUMN ,false},
	"commit":{TOKEN_SQL_COMMIT ,false},
	"conflict":{TOKEN_SQL_CONFLICT ,false},
	"constraint":{TOKEN_SQL_CONSTRAINT ,false},
	"create":{TOKEN_SQL_CREATE ,false},
	"cross":{TOKEN_SQL_CROSS ,false},
	"current":{TOKEN_SQL_CURRENT ,false},
	"current_date":{TOKEN_SQL_CURRENT_DATE ,false},
	"current_time":{TOKEN_SQL_CURRENT_TIME ,false},
	"current_timestamp":{TOKEN_SQL_CURRENT_TIMESTAMP ,false},
	"database":{TOKEN_SQL_DATABASE ,false},
	"default":{TOKEN_SQL_DEFAULT ,false},
	"deferrable":{TOKEN_SQL_DEFERRABLE ,false},
	"deferred":{TOKEN_SQL_DEFERRED ,false},
	"delete":{TOKEN_SQL_DELETE ,false},
	"desc":{TOKEN_SQL_DESC ,false},
	"detach":{TOKEN_SQL_DETACH ,false},
	"distinct":{TOKEN_SQL_DISTINCT ,false},
	"do":{TOKEN_SQL_DO ,false},
	"drop":{TOKEN_SQL_DROP ,false},
	"each":{TOKEN_SQL_EACH ,false},
	"else":{TOKEN_SQL_ELSE ,false},
	"end":{TOKEN_SQL_END ,false},
	"escape":{TOKEN_SQL_ESCAPE ,false},
	"except":{TOKEN_SQL_EXCEPT ,false},
	"exclude":{TOKEN_SQL_EXCLUDE ,false},
	"exclusive":{TOKEN_SQL_EXCLUSIVE ,false},
	"exists":{TOKEN_SQL_EXISTS ,false},
	"explain":{TOKEN_SQL_EXPLAIN ,false},
	"fail":{TOKEN_SQL_FAIL ,false},
	"filter":{TOKEN_SQL_FILTER ,false},
	"first":{TOKEN_SQL_FIRST ,false},
	"following":{TOKEN_SQL_FOLLOWING ,false},
	"for":{TOKEN_SQL_FOR ,false},
	"foreign":{TOKEN_SQL_FOREIGN ,false},
	"from":{TOKEN_SQL_FROM ,false},
	"full":{TOKEN_SQL_FULL ,false},
	"generated":{TOKEN_SQL_GENERATED ,false},
	"glob":{TOKEN_SQL_GLOB ,false},
	"group":{TOKEN_SQL_GROUP ,false},
	"groups":{TOKEN_SQL_GROUPS ,false},
	"having":{TOKEN_SQL_HAVING ,false},
	"if":{TOKEN_SQL_IF ,false},
	"ignore":{TOKEN_SQL_IGNORE ,false},
	"immediate":{TOKEN_SQL_IMMEDIATE ,false},
	"in":{TOKEN_SQL_IN ,false},
	"index":{TOKEN_SQL_INDEX ,false},
	"indexed":{TOKEN_SQL_INDEXED ,false},
	"initially":{TOKEN_SQL_INITIALLY ,false},
	"inner":{TOKEN_SQL_INNER ,false},
	"insert":{TOKEN_SQL_INSERT ,false},
	"instead":{TOKEN_SQL_INSTEAD ,false},
	"intersect":{TOKEN_SQL_INTERSECT ,false},
	"into":{TOKEN_SQL_INTO ,false},
	"is":{TOKEN_SQL_IS ,false},
	"isnull":{TOKEN_SQL_ISNULL ,false},
	"join":{TOKEN_SQL_JOIN ,false},
	"key":{TOKEN_SQL_KEY ,false},
	"last":{TOKEN_SQL_LAST ,false},
	"left":{TOKEN_SQL_LEFT ,false},
	"like":{TOKEN_SQL_LIKE ,false},
	"limit":{TOKEN_SQL_LIMIT ,false},
	"match":{TOKEN_SQL_MATCH ,false},
	"natural":{TOKEN_SQL_NATURAL ,false},
	"no":{TOKEN_SQL_NO ,false},
	"not":{TOKEN_SQL_NOT ,false},
	"nothing":{TOKEN_SQL_NOTHING ,false},
	"notnull":{TOKEN_SQL_NOTNULL ,false},
	"null":{TOKEN_SQL_NULL ,false},
	"nulls":{TOKEN_SQL_NULLS ,false},
	"of":{TOKEN_SQL_OF ,false},
	"offset":{TOKEN_SQL_OFFSET ,false},
	"on":{TOKEN_SQL_ON ,false},
	"or":{TOKEN_SQL_OR ,false},
	"order":{TOKEN_SQL_ORDER ,false},
	"others":{TOKEN_SQL_OTHERS ,false},
	"outer":{TOKEN_SQL_OUTER ,false},
	"over":{TOKEN_SQL_OVER ,false},
	"partition":{TOKEN_SQL_PARTITION ,false},
	"plan":{TOKEN_SQL_PLAN ,false},
	"pragma":{TOKEN_SQL_PRAGMA ,false},
	"preceding":{TOKEN_SQL_PRECEDING ,false},
	"primary":{TOKEN_SQL_PRIMARY ,false},
	"query":{TOKEN_SQL_QUERY ,false},
	"raise":{TOKEN_SQL_RAISE ,false},
	"range":{TOKEN_SQL_RANGE ,false},
	"recursive":{TOKEN_SQL_RECURSIVE ,false},
	"references":{TOKEN_SQL_REFERENCES ,false},
	"regexp":{TOKEN_SQL_REGEXP ,false},
	"reindex":{TOKEN_SQL_REINDEX ,false},
	"release":{TOKEN_SQL_RELEASE ,false},
	"rename":{TOKEN_SQL_RENAME ,false},
	"replace":{TOKEN_SQL_REPLACE ,false},
	"restrict":{TOKEN_SQL_RESTRICT ,false},
	"right":{TOKEN_SQL_RIGHT ,false},
	"rollback":{TOKEN_SQL_ROLLBACK ,false},
	"row":{TOKEN_SQL_ROW ,false},
	"rows":{TOKEN_SQL_ROWS ,false},
	"savepoint":{TOKEN_SQL_SAVEPOINT ,false},
	"select":{TOKEN_SQL_SELECT ,false},
	"set":{TOKEN_SQL_SET ,false},
	"table":{TOKEN_SQL_TABLE ,false},
	"temp":{TOKEN_SQL_TEMP ,false},
	"temporary":{TOKEN_SQL_TEMPORARY ,false},
	"then":{TOKEN_SQL_THEN ,false},
	"ties":{TOKEN_SQL_TIES ,false},
	"to":{TOKEN_SQL_TO ,false},
	"transaction":{TOKEN_SQL_TRANSACTION ,false},
	"trigger":{TOKEN_SQL_TRIGGER ,false},
	"unbounded":{TOKEN_SQL_UNBOUNDED ,false},
	"union":{TOKEN_SQL_UNION ,false},
	"unique":{TOKEN_SQL_UNIQUE ,false},
	"update":{TOKEN_SQL_UPDATE ,false},
	"using":{TOKEN_SQL_USING ,false},
	"vacuum":{TOKEN_SQL_VACUUM ,false},
	"values":{TOKEN_SQL_VALUES ,false},
	"view":{TOKEN_SQL_VIEW ,false},
	"virtual":{TOKEN_SQL_VIRTUAL ,false},
	"when":{TOKEN_SQL_WHEN ,false},
	"where":{TOKEN_SQL_WHERE ,false},
	"window":{TOKEN_SQL_WINDOW ,false},
	"with":{TOKEN_SQL_WITH ,false},
	"without":{TOKEN_SQL_WITHOUT ,false},
}


type Precedence byte

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_INCR
	PREC_FACTOR // * /
	PREC_UNARY  // ! -
	PREC_CALL   // . ()
	PREC_ARRAY // . @[]
	PREC_INDEX
	PREC_LIST  //. 	@{}
	PREC_PRIMARY
)
