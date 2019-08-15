%{
package parser
%}

%token BYTE
%token ELSE
%token FALSE
%token IF
%token INT
%token PROC
%token REFERENCE
%token RETURN
%token WHILE
%token TRUE
%token IDENT
%token INT_CONST
%token CHAR_LIT
%token STR_LIT
%token EQ
%token NE
%token LE
%token GE

/* Deal with dangling else */
%nonassoc ')'
%nonassoc ELSE
%left '|'
%left '&'
%nonassoc EQ NE '<' '>' LE GE
%left '+' '-'
%left '*' '/' '%'
%nonassoc SIGN

%union {
    id string
	i  int
	c  byte
	s  string
}

%%

program:
	func_def
;

func_def:
	IDENT '(' ')' ':' r_type local_def_list compound_stmt
|	IDENT '(' fpar_list ')' ':' r_type local_def_list compound_stmt
;

fpar_list:
	fpar_def
|	fpar_list ',' fpar_def
;

fpar_def:
	IDENT ':' type
|	IDENT ':' REFERENCE type
;

data_type:
	INT
|	BYTE
;

type:
	data_type
|	data_type '[' ']'
;

r_type:
	data_type
|	PROC
;

local_def_list:
	/* empty */
|	local_def_list local_def
;

local_def:
	func_def
|	var_def
;

var_def:
	IDENT ':' data_type ';'
|	IDENT ':' data_type '[' INT_CONST ']' ';'
;

stmt:
	';'
|	l_value '=' expr ';'
|	compound_stmt
|	func_call ';'
|	IF '(' cond ')' stmt
|	IF '(' cond ')' stmt ELSE stmt
|	WHILE '(' cond ')' stmt
|	RETURN ';'
|	RETURN expr ';'
;

compound_stmt:
	'{' stmt_list '}'
;

stmt_list:
	/* empty */
|	stmt_list stmt
;

func_call:
	IDENT '(' ')'
|	IDENT '(' expr_list ')'
;

expr_list:
	expr
|	expr_list ',' expr
;

expr:
	INT_CONST
|	CHAR_LIT
|	l_value
|	'(' expr ')'
|	func_call
|	'+' expr %prec SIGN
|	'-' expr %prec SIGN
|	expr '+' expr
|	expr '-' expr
|	expr '*' expr
|	expr '/' expr
|	expr '%' expr
;

l_value:
	IDENT
|	IDENT '[' expr ']'
|	STR_LIT
;

cond:
	TRUE
|	FALSE
|	'(' cond ')'
|	'!' cond %prec SIGN
|	expr EQ expr
|	expr NE expr
|	expr '<' expr
|	expr '>' expr
|	expr LE expr
|	expr GE expr
|	cond '&' cond
|	cond '|' cond
;

%%
