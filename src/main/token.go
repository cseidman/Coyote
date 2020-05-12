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
	// 70
	TOKEN_RIGHT
	TOKEN_CROSSJOIN
	TOKEN_WHERE
	TOKEN_ALL
	TOKEN_ORDER
	TOKEN_GROUP
	TOKEN_BY
	TOKEN_INTO
	TOKEN_VALUES
	TOKEN_AS
	// 80
	TOKEN_ON
	TOKEN_BROWSE
	TOKEN_CREATE
	TOKEN_TABLE
	TOKEN_COLUMN
	TOKEN_ROW
	TOKEN_VIEW
	TOKEN_HAVING
	TOKEN_DISTINCT
	TOKEN_TOP
	//90
	TOKEN_STEP
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
	"++":       {TOKEN_PLUS_PLUS, true},
	"--":       {TOKEN_MINUS_MINUS, true},
	"@[":       {TOKEN_ARRAY, true},
	"@{":       {TOKEN_LIST, true},
	"and":      {TOKEN_AND, true},
	"class":    {TOKEN_CLASS, true},
	"else":     {TOKEN_ELSE, true},
	"false":    {TOKEN_FALSE, true},
	"for":      {TOKEN_FOR, true},
	"func": 	{TOKEN_FUNC, true},
	"select":   {TOKEN_SELECT, true},
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
	"_array":   {TOKEN_TYPE_ARRAY, true},
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
	"list": {TOKEN_LIST_TYPE, true},
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
	PREC_INDEX
	PREC_ARRAY // . @[]
	PREC_LIST  //. 	@{}
	PREC_PRIMARY
)
