%{
package parser

import "github.com/foxeng/alanc/ast"

// _ast is the AST constructed by the parser.
// TODO OPT: Avoid this global. How?
var _ast *ast.Ast
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
	id ast.ID
	ast ast.Ast
	fdef ast.FuncDef
	pdefs []ast.ParDef
	pdef ast.ParDef
	dt ast.DataType
	rtype *ast.DataType
	ldefs []ast.LocalDef
	ldef ast.LocalDef
	vdef ast.VarDef
	stmt ast.Stmt
	cstmt ast.CompStmt
	stmts []ast.Stmt
	fcall ast.FuncCall
	exprs []ast.Expr
	expr ast.Expr
	lval ast.LVal
	cond ast.Cond
	iconst ast.IntConstExpr
	cconst ast.CharConstExpr
	strlit ast.StrLitExpr
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
		_ast = &ast.Ast{
			Program: $1,
		}
	}
;

func_def:
	IDENT '(' fpar_list ')' ':' r_type local_def_list compound_stmt
	{
		$$ = ast.FuncDef{
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
		$$ = []ast.ParDef{}
	}
|	fpar_def
	{
		$$ = []ast.ParDef{$1}
	}
|	fpar_list ',' fpar_def
	{
		$$ = append($1, $3)
	}
;

fpar_def:
	IDENT ':' data_type
	{
		$$ = ast.ParDef{
			VarDef: &ast.PrimVarDef{
				ID: $1,
				DataType: $3,
			},
		}
	}
|	IDENT ':' REFERENCE data_type
	{
		$$ = ast.ParDef{
			VarDef: &ast.PrimVarDef{
				ID: $1,
				DataType: $4,
			},
			IsRef: true,
		}
	}
|	IDENT ':' REFERENCE data_type '[' ']'
	{
		$$ = ast.ParDef{
			VarDef: &ast.ArrayDef{
				PrimVarDef: ast.PrimVarDef{
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
		$$ = ast.DataTypeInt
	}
|	BYTE
	{
		$$ = ast.DataTypeByte
	}
;

r_type:
	data_type
	{
		$$ = &$1
	}
|	PROC
	{
		$$ = nil
	}
;

local_def_list:
	/* empty */
	{
		$$ = []ast.LocalDef{}
	}
|	local_def_list local_def
	{
		$$ = append($1, $2)
	}
;

local_def:
	func_def
	{
		$$ = &$1
	}
|	var_def
	{
		$$ = $1
	}
;

var_def:
	IDENT ':' data_type ';'
	{
		$$ = &ast.PrimVarDef{
			ID: $1,
			DataType: $3,
		}
	}
|	IDENT ':' data_type '[' INT_CONST ']' ';'
	{
		$$ = &ast.ArrayDef{
			PrimVarDef: ast.PrimVarDef{
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
		$$ = &ast.CompStmt{
			Stmts: []ast.Stmt{},
		}
	}
|	l_value '=' expr ';'
	{
		$$ = &ast.AssignStmt{
			Left: $1,
			Right: $3,
		}
	}
|	compound_stmt
	{
		$$ = &$1
	}
|	func_call ';'
	{
		$$ = &ast.FuncCallStmt{
			FuncCall: $1,
		}
	}
|	IF '(' cond ')' stmt
	{
		$$ = &ast.IfStmt{
			Cond: $3,
			Stmt: $5,
		}
	}
|	IF '(' cond ')' stmt ELSE stmt
	{
		$$ = &ast.IfElseStmt{
			Cond: $3,
			Stmt1: $5,
			Stmt2: $7,
		}
	}
|	WHILE '(' cond ')' stmt
	{
		$$ = &ast.WhileStmt{
			Cond: $3,
			Stmt: $5,
		}
	}
|	RETURN ';'
	{
		$$ = &ast.ReturnStmt{
			Expr: nil,
		}
	}
|	RETURN expr ';'
	{
		$$ = &ast.ReturnStmt{
			Expr: $2,
		}
	}
;

compound_stmt:
	'{' stmt_list '}'
	{
		$$ = ast.CompStmt{
			Stmts: $2,
		}
	}
;

stmt_list:
	/* empty */
	{
		$$ = []ast.Stmt{}
	}
|	stmt_list stmt
	{
		$$ = append($1, $2)
	}
;

func_call:
	IDENT '(' ')'
	{
		$$ = ast.FuncCall{
			ID: $1,
			Args: []ast.Expr{},
		}
	}
|	IDENT '(' expr_list ')'
	{
		$$ = ast.FuncCall{
			ID: $1,
			Args: $3,
		}
	}
;

expr_list:
	expr
	{
		$$ = []ast.Expr{$1}
	}
|	expr_list ',' expr
	{
		$$ = append($1, $3)
	}
;

expr:
	INT_CONST
	{
		$$ = &$1
	}
|	CHAR_LIT
	{
		$$ = &$1
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
		$$ = &ast.FuncCallExpr{
			FuncCall: $1,
		}
	}
|	'+' expr %prec SIGN
	{
		$$ = &ast.UnArithExpr{
			Sign: ast.SignPlus,
			Expr: $2,
		}
	}
|	'-' expr %prec SIGN
	{
		$$ = &ast.UnArithExpr{
			Sign: ast.SignMinus,
			Expr: $2,
		}
	}
|	expr '+' expr
	{
		$$ = &ast.BinArithExpr{
			Left: $1,
			Op: ast.ArithOpPlus,
			Right: $3,
		}
	}
|	expr '-' expr
	{
		$$ = &ast.BinArithExpr{
			Left: $1,
			Op: ast.ArithOpMinus,
			Right: $3,
		}
	}
|	expr '*' expr
	{
		$$ = &ast.BinArithExpr{
			Left: $1,
			Op: ast.ArithOpMult,
			Right: $3,
		}
	}
|	expr '/' expr
	{
		$$ = &ast.BinArithExpr{
			Left: $1,
			Op: ast.ArithOpDiv,
			Right: $3,
		}
	}
|	expr '%' expr
	{
		$$ = &ast.BinArithExpr{
			Left: $1,
			Op: ast.ArithOpMod,
			Right: $3,
		}
	}
;

l_value:
	IDENT
	{
		$$ = &ast.VarRef{
			ID: $1,
		}
	}
|	IDENT '[' expr ']'
	{
		$$ = &ast.ArrayElem{
			ID: $1,
			Index: $3,
		}
	}
|	STR_LIT
	{
		$$ = &$1
	}
;

cond:
	TRUE
	{
		$$ = &ast.ConstCond{
			Val: true,
		}
	}
|	FALSE
	{
		$$ = &ast.ConstCond{
			Val: false,
		}
	}
|	'(' cond ')'
	{
		$$ = $2
	}
|	'!' cond %prec SIGN
	{
		$$ = &ast.UnCond{
			Cond: $2,
		}
	}
|	expr EQ expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpEQ,
			Right: $3,
		}
	}
|	expr NE expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpNE,
			Right: $3,
		}
	}
|	expr '<' expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpLT,
			Right: $3,
		}
	}
|	expr '>' expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpGT,
			Right: $3,
		}
	}
|	expr LE expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpLE,
			Right: $3,
		}
	}
|	expr GE expr
	{
		$$ = &ast.CompCond{
			Left: $1,
			Op: ast.CompOpGE,
			Right: $3,
		}
	}
|	cond '&' cond
	{
		$$ = &ast.BinCond{
			Left: $1,
			Op: ast.LogOpAnd,
			Right: $3,
		}
	}
|	cond '|' cond
	{
		$$ = &ast.BinCond{
			Left: $1,
			Op: ast.LogOpOr,
			Right: $3,
		}
	}
;

%%
