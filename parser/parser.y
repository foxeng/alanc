%{
package parser

import "github.com/foxeng/alanc/semantic"

// ast is the AST constructed by the parser.
// TODO OPT: Avoid this global. How?
var ast *semantic.Ast
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
	id semantic.ID
	ast semantic.Ast
	fdef *semantic.FuncDef
	pdefs []semantic.ParDef
	pdef semantic.ParDef
	dt semantic.DataType
	rtype *semantic.DataType
	ldefs []semantic.LocalDef
	ldef semantic.LocalDef
	vdef semantic.VarDef
	stmt semantic.Stmt
	cstmt semantic.CompStmt
	stmts []semantic.Stmt
	fcall semantic.FuncCall
	exprs []semantic.Expr
	expr semantic.Expr
	lval semantic.LVal
	cond semantic.Cond
	iconst semantic.IntConstExpr
	cconst semantic.CharConstExpr
	strlit semantic.StrLitExpr
}

%type <yys> BYTE
			ELSE
			FALSE
			IF
			INT
			PROC
			REFERENCE
			RETURN
			WHILE
			TRUE
			EQ
			NE
			LE
			GE
%type <id> IDENT
%type <ast> program
%type <fdef> func_def
%type <pdefs> fpar_list
%type <rtype> r_type
%type <ldefs> local_def_list
%type <cstmt> compound_stmt
%type <pdef> fpar_def
%type <dt> data_type
%type <ldef> local_def
%type <vdef> var_def
%type <stmt> stmt
%type <lval> l_value
%type <expr> expr
%type <fcall> func_call
%type <cond> cond
%type <stmts> stmt_list
%type <exprs> expr_list
%type <iconst> INT_CONST
%type <cconst> CHAR_LIT
%type <strlit> STR_LIT

%%

program:
	func_def
	{
		ast = &semantic.Ast{
			Program: $1,
		}
	}
;

func_def:
	IDENT '(' fpar_list ')' ':' r_type local_def_list compound_stmt
	{
		$$ = &semantic.FuncDef{
			ID: $1,
			Parameters: $3,
			RType: $6,
			LDefs: $7,
			CompStmt: $8,
		}
	}
;

fpar_list:
	/* empty */
	{
		$$ = []semantic.ParDef{}
	}
|	fpar_def
	{
		$$ = []semantic.ParDef{$1}
	}
|	fpar_list ',' fpar_def
	{
		$$ = append($1, $3)
	}
;

fpar_def:
	IDENT ':' data_type
	{
		$$ = semantic.ParDef{
			VarDef: &semantic.PrimVarDef{
				ID: $1,
				DataType: $3,
			},
		}
	}
|	IDENT ':' REFERENCE data_type
	{
		$$ = semantic.ParDef{
			VarDef: &semantic.PrimVarDef{
				ID: $1,
				DataType: $4,
			},
			IsRef: true,
		}
	}
|	IDENT ':' REFERENCE data_type '[' ']'
	{
		$$ = semantic.ParDef{
			VarDef: &semantic.ArrayDef{
				PrimVarDef: semantic.PrimVarDef{
					ID: $1,
					DataType: $4,
				},
				// TODO OPT: Set size to some special value indicating unknown?
			},
			IsRef: true,
		}
	}
;

data_type:
	INT
	{
		$$ = semantic.DataTypeInt
	}
|	BYTE
	{
		$$ = semantic.DataTypeByte
	}
;

r_type:
	data_type
	{
		// Make sure to return a pointer to a _copy_ of $1, because the underlying $1 is reused.
		dt := $1
		$$ = &dt
	}
|	PROC
	{
		$$ = nil
	}
;

local_def_list:
	/* empty */
	{
		$$ = []semantic.LocalDef{}
	}
|	local_def_list local_def
	{
		$$ = append($1, $2)
	}
;

local_def:
	func_def
	{
		$$ = $1
	}
|	var_def
	{
		$$ = $1
	}
;

var_def:
	IDENT ':' data_type ';'
	{
		$$ = &semantic.PrimVarDef{
			ID: $1,
			DataType: $3,
		}
	}
|	IDENT ':' data_type '[' INT_CONST ']' ';'
	{
		$$ = &semantic.ArrayDef{
			PrimVarDef: semantic.PrimVarDef{
				ID: $1,
				DataType: $3,
			},
			Size: $5,
		}
	}
;

stmt:
	';'
	{
		$$ = &semantic.CompStmt{
			Stmts: []semantic.Stmt{},
		}
	}
|	l_value '=' expr ';'
	{
		$$ = &semantic.AssignStmt{
			Left: $1,
			Right: $3,
		}
	}
|	compound_stmt
	{
		// Make sure to return a pointer to a _copy_ of $1, because the underlying $1 is reused.
		cs := $1
		$$ = &cs
	}
|	func_call ';'
	{
		$$ = &semantic.FuncCallStmt{
			FuncCall: $1,
		}
	}
|	IF '(' cond ')' stmt
	{
		$$ = &semantic.IfStmt{
			Cond: $3,
			Stmt: $5,
		}
	}
|	IF '(' cond ')' stmt ELSE stmt
	{
		$$ = &semantic.IfElseStmt{
			Cond: $3,
			Stmt1: $5,
			Stmt2: $7,
		}
	}
|	WHILE '(' cond ')' stmt
	{
		$$ = &semantic.WhileStmt{
			Cond: $3,
			Stmt: $5,
		}
	}
|	RETURN ';'
	{
		$$ = &semantic.ReturnStmt{
			Expr: nil,
		}
	}
|	RETURN expr ';'
	{
		$$ = &semantic.ReturnStmt{
			Expr: $2,
		}
	}
;

compound_stmt:
	'{' stmt_list '}'
	{
		$$ = semantic.CompStmt{
			Stmts: $2,
		}
	}
;

stmt_list:
	/* empty */
	{
		$$ = []semantic.Stmt{}
	}
|	stmt_list stmt
	{
		$$ = append($1, $2)
	}
;

func_call:
	IDENT '(' ')'
	{
		$$ = semantic.FuncCall{
			ID: $1,
			Args: []semantic.Expr{},
		}
	}
|	IDENT '(' expr_list ')'
	{
		$$ = semantic.FuncCall{
			ID: $1,
			Args: $3,
		}
	}
;

expr_list:
	expr
	{
		$$ = []semantic.Expr{$1}
	}
|	expr_list ',' expr
	{
		$$ = append($1, $3)
	}
;

expr:
	INT_CONST
	{
		// Make sure to return a pointer to a _copy_ of $1, because the underlying $1 is reused.
		i := $1
		$$ = &i
	}
|	CHAR_LIT
	{
		// Make sure to return a pointer to a _copy_ of $1, because the underlying $1 is reused.
		c := $1
		$$ = &c
	}
|	l_value
	{
		$$ = $1
	}
|	'(' expr ')'
	{
		$$ = $2
	}
|	func_call
	{
		$$ = &semantic.FuncCallExpr{
			FuncCall: $1,
		}
	}
|	'+' expr %prec SIGN
	{
		$$ = &semantic.UnArithExpr{
			Sign: semantic.SignPlus,
			Expr: $2,
		}
	}
|	'-' expr %prec SIGN
	{
		$$ = &semantic.UnArithExpr{
			Sign: semantic.SignMinus,
			Expr: $2,
		}
	}
|	expr '+' expr
	{
		$$ = &semantic.BinArithExpr{
			Left: $1,
			Op: semantic.ArithOpPlus,
			Right: $3,
		}
	}
|	expr '-' expr
	{
		$$ = &semantic.BinArithExpr{
			Left: $1,
			Op: semantic.ArithOpMinus,
			Right: $3,
		}
	}
|	expr '*' expr
	{
		$$ = &semantic.BinArithExpr{
			Left: $1,
			Op: semantic.ArithOpMult,
			Right: $3,
		}
	}
|	expr '/' expr
	{
		$$ = &semantic.BinArithExpr{
			Left: $1,
			Op: semantic.ArithOpDiv,
			Right: $3,
		}
	}
|	expr '%' expr
	{
		$$ = &semantic.BinArithExpr{
			Left: $1,
			Op: semantic.ArithOpMod,
			Right: $3,
		}
	}
;

l_value:
	IDENT
	{
		$$ = &semantic.VarRef{
			ID: $1,
		}
	}
|	IDENT '[' expr ']'
	{
		$$ = &semantic.ArrayElem{
			ID: $1,
			Index: $3,
		}
	}
|	STR_LIT
	{
		// Make sure to return a pointer to a _copy_ of $1, because the underlying $1 is reused.
		s := $1
		$$ = &s
	}
;

cond:
	TRUE
	{
		$$ = &semantic.ConstCond{
			Val: true,
		}
	}
|	FALSE
	{
		$$ = &semantic.ConstCond{
			Val: false,
		}
	}
|	'(' cond ')'
	{
		$$ = $2
	}
|	'!' cond %prec SIGN
	{
		$$ = &semantic.UnCond{
			Cond: $2,
		}
	}
|	expr EQ expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpEQ,
			Right: $3,
		}
	}
|	expr NE expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpNE,
			Right: $3,
		}
	}
|	expr '<' expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpLT,
			Right: $3,
		}
	}
|	expr '>' expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpGT,
			Right: $3,
		}
	}
|	expr LE expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpLE,
			Right: $3,
		}
	}
|	expr GE expr
	{
		$$ = &semantic.CompCond{
			Left: $1,
			Op: semantic.CompOpGE,
			Right: $3,
		}
	}
|	cond '&' cond
	{
		$$ = &semantic.BinCond{
			Left: $1,
			Op: semantic.LogOpAnd,
			Right: $3,
		}
	}
|	cond '|' cond
	{
		$$ = &semantic.BinCond{
			Left: $1,
			Op: semantic.LogOpOr,
			Right: $3,
		}
	}
;

%%
