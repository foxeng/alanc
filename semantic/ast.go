// Package semantic contains the semantic analysis for Alan: AST definition, symbol table, type
// checking.
package semantic

/*
AST nodes:
function definition: id, [parameter list], return type, [local definitions list], statement
variable definition: id, type
array definition: id, type, size

empty
seq: statement, statement
assignment: l-value, expression
function call: id, expression list
if: condition, statement
if else: condition, statement, statement
while: condition, statement
return: [expression]

int constant: int
char constant: char
function call: id, expression list
string literal: string
unary arithmetic: op (+, -), expression
binary arithmetic: op (+, -, *, /, %), expression, expression

variable reference: id
array element: id, expression

bool constant: bool
comparison: op (==, !=, <, >, <=, >=), expression, expression
unary condition: condition
binary condition: op (&, |), condition, condition
*/

const (
	// DataTypeInt is the "int" data type.
	DataTypeInt DataType = iota
	// DataTypeByte is the "byte" data type.
	DataTypeByte
)

const (
	// SignPlus is the '+' sign.
	SignPlus Sign = '+'
	// SignMinus is the '-' sign.
	SignMinus Sign = '-'
)

const (
	// ArithOpPlus is the '+' arithmetic operator.
	ArithOpPlus ArithOp = '+'
	// ArithOpMinus is the '-' arithmetic operator.
	ArithOpMinus ArithOp = '-'
	// ArithOpMult is the '*' arithmetic operator.
	ArithOpMult ArithOp = '*'
	// ArithOpDiv is the '/' arithmetic operator.
	ArithOpDiv ArithOp = '/'
	// ArithOpMod is the '%' arithmetic operator.
	ArithOpMod ArithOp = '%'
)

const (
	// CompOpEQ is the "==" comparison operator.
	CompOpEQ CompOp = "=="
	// CompOpNE is the "!=" comparison operator.
	CompOpNE CompOp = "!="
	// CompOpLT is the "<" comparison operator.
	CompOpLT CompOp = "<"
	// CompOpGT is the ">" comparison operator.
	CompOpGT CompOp = ">"
	// CompOpLE is the "<=" comparison operator.
	CompOpLE CompOp = "<="
	// CompOpGE is the ">=" comparison operator.
	CompOpGE CompOp = ">="
)

const (
	// LogOpAnd is the '&' logical operator.
	LogOpAnd LogOp = '&'
	// LogOpOr is the '|' logical operator.
	LogOpOr LogOp = '|'
)

// ID is an identifier.
type ID string

// DataType is a (primitive) data type (i.e. "int" or "byte").
type DataType int

// Sign is an arithmetic sign (i.e. '+' or '-').
type Sign rune

// ArithOp is an arithmetic operator (i.e. '+', '-', '*', '/' or '%').
type ArithOp rune

//CompOp is a comparison operator (i.e. "==", "!=", "<", ">", "<=" or ">=").
type CompOp string

// LogOp is a logical operator (i.e. '&' or '|').
type LogOp rune

// Ast is a whole abstract syntax tree.
type Ast struct {
	Program *FuncDef
}

// Node is a single Node of an AST.
type Node interface {
	isNode()
}

// LocalDef is a local definition (discriminated union of FuncDef and VarDef).
type LocalDef interface {
	Node
	isLocalDef()
	Id() ID
}

// FuncDef is a function definition.
type FuncDef struct {
	// ID is the function's identifier.
	ID
	// Parameters are the function's parameters.
	Parameters []ParDef
	// RType is the function's return type (nil if function is a proc).
	RType *DataType
	// LDefs are the function's local definitions.
	LDefs []LocalDef
	// CompStmt is the function's body.
	CompStmt
}

// TODO OPT: Define the methods on value instead of pointer receivers?

func (*FuncDef) isNode() {}

func (*FuncDef) isLocalDef() {}

// Id returns the function's identifier.
func (fd *FuncDef) Id() ID {
	return fd.ID
}

// ParDef is a function parameter's definition.
type ParDef struct {
	// VarDef is the underlying variable definition (if an array, size is ignored).
	VarDef
	// IsRef denotes whether the parameter is passed by reference.
	IsRef bool // TODO OPT: Avoid this by making separate ValParDef and RefParDef?
}

func (*ParDef) isNode() {}

func (*ParDef) isLocalDef() {}

// Id returns the parameter's identifier.
func (pd *ParDef) Id() ID {
	return pd.VarDef.Id()
}

// VarDef is the discriminated union of PrimVarDef and ArrayDef.
type VarDef interface {
	LocalDef
	isVarDef()
}

// PrimVarDef is a primitive variable definition.
type PrimVarDef struct {
	// ID is the variable's identifier.
	ID
	// DataType is the variable's (primitive) data type.
	DataType
}

func (*PrimVarDef) isNode() {}

func (*PrimVarDef) isLocalDef() {}

// Id returns the primitive variable's identifier.
func (pvd *PrimVarDef) Id() ID {
	return pvd.ID
}

func (*PrimVarDef) isVarDef() {}

// ArrayDef is an array definition.
type ArrayDef struct {
	// PrimVarDef is the underlying variable definition.
	PrimVarDef
	// Size is the size of the array
	Size IntConstExpr
}

func (*ArrayDef) isNode() {}

func (*ArrayDef) isLocalDef() {}

// Id returns the array's identifier.
func (ad *ArrayDef) Id() ID {
	return ad.ID
}

func (*ArrayDef) isVarDef() {}

// Stmt is a statement.
type Stmt interface {
	Node
	isStmt()
}

// CompStmt is a compound statement.
type CompStmt struct {
	// Stmts are the statement's constituents.
	Stmts []Stmt
}

func (*CompStmt) isNode() {}

func (*CompStmt) isStmt() {}

// AssignStmt is an assignment statement.
type AssignStmt struct {
	// Left is the left-hand side of the assignment.
	Left LVal
	// Right is the right-hand side of the assignment.
	Right Expr
}

func (*AssignStmt) isNode() {}

func (*AssignStmt) isStmt() {}

// FuncCall is a function call.
type FuncCall struct {
	// ID is the function's identifier.
	ID
	// Args are the call's arguments.
	Args []Expr
}

func (*FuncCall) isNode() {}

// FuncCallStmt is a function call statement.
type FuncCallStmt struct {
	// FuncCall is the underlying function call.
	FuncCall
}

func (*FuncCallStmt) isNode() {}

func (*FuncCallStmt) isStmt() {}

// IfStmt is an if statement.
type IfStmt struct {
	// Cond is the if statement's condition.
	Cond
	// Stmt is the if statement's body.
	Stmt
}

func (*IfStmt) isNode() {}

func (*IfStmt) isStmt() {}

// TODO: Merge with IfStmt?

// IfElseStmt is an if-else statement.
type IfElseStmt struct {
	// Cond is the if statement's condition.
	Cond
	// Stmt1 is the if clause's body.
	Stmt1 Stmt
	// Stmt2 is the else clause's body.
	Stmt2 Stmt
}

func (*IfElseStmt) isNode() {}

func (*IfElseStmt) isStmt() {}

// WhileStmt is a while statement.
type WhileStmt struct {
	// Cond is the while statement's condition.
	Cond
	// Stmt is the while statement's body.
	Stmt
}

func (*WhileStmt) isNode() {}

func (*WhileStmt) isStmt() {}

// ReturnStmt is a return statement.
type ReturnStmt struct {
	// Expr is the return expression (nil if nothing is returned).
	Expr
}

func (*ReturnStmt) isNode() {}

func (*ReturnStmt) isStmt() {}

// Expr is an expression.
type Expr interface {
	Node
	isExpr()
}

// IntConstExpr is an integer constant expression.
type IntConstExpr struct {
	Val int // TODO OPT: Use fixed width?
}

func (*IntConstExpr) isNode() {}

func (*IntConstExpr) isExpr() {}

// CharConstExpr is a character constant expression.
type CharConstExpr struct {
	Val rune
}

func (*CharConstExpr) isNode() {}

func (*CharConstExpr) isExpr() {}

// LVal is an l-value.
type LVal interface {
	Expr
	isLVal()
}

// VarRef is a variable reference.
type VarRef struct {
	// ID is the variable's identifier.
	ID
}

func (*VarRef) isNode() {}

func (*VarRef) isExpr() {}

func (*VarRef) isLVal() {}

// TODO: Merge with VarRef?

// ArrayElem is an array element.
type ArrayElem struct {
	// ID is the array's identifier.
	ID
	// Index is the element's index.
	Index Expr
}

func (*ArrayElem) isNode() {}

func (*ArrayElem) isExpr() {}

func (*ArrayElem) isLVal() {}

// StrLitExpr is a string literal expression.
type StrLitExpr struct {
	// Val is the underlying string literal.
	Val string
}

func (*StrLitExpr) isNode() {}

func (*StrLitExpr) isExpr() {}

func (*StrLitExpr) isLVal() {}

// FuncCallExpr is a function call expression.
type FuncCallExpr struct {
	// FuncCall is the underlying function call.
	FuncCall
}

func (*FuncCallExpr) isNode() {}

func (*FuncCallExpr) isExpr() {}

// UnArithExpr is an unary arithmetic expression.
type UnArithExpr struct {
	// Sign is the expression's sign.
	Sign
	// Expr is the underlying (arithmetic) expression.
	Expr
}

func (*UnArithExpr) isNode() {}

func (*UnArithExpr) isExpr() {}

// BinArithExpr is a binary arithmetic expression.
type BinArithExpr struct {
	// Left is the left-hand side of the expression.
	Left Expr
	// Op is the expression's operator.
	Op ArithOp
	// Right is the right-hand side of the expression.
	Right Expr
}

func (*BinArithExpr) isNode() {}

func (*BinArithExpr) isExpr() {}

// Cond is a condition.
type Cond interface {
	Expr
	isCond()
}

// ConstCond is a constant condition.
type ConstCond struct {
	// Val is the underlying constant (i.e. true or false)
	Val bool
}

func (*ConstCond) isNode() {}

func (*ConstCond) isExpr() {}

func (*ConstCond) isCond() {}

// UnCond is an unary condition (negation).
type UnCond struct {
	// Cond is the underlying condition.
	Cond
}

func (*UnCond) isNode() {}

func (*UnCond) isExpr() {}

func (*UnCond) isCond() {}

// CompCond is a comparison condition.
type CompCond struct {
	// Left is the left-hand side of the comparison.
	Left Expr
	// Op is the comparison operator.
	Op CompOp
	// Right is the right-hand side of the comparison.
	Right Expr
}

func (*CompCond) isNode() {}

func (*CompCond) isExpr() {}

func (*CompCond) isCond() {}

// BinCond is a binary logical condition.
type BinCond struct {
	// Left is the left-hand side of the condition.
	Left Cond
	// Op is the logical operator.
	Op LogOp
	// Right is the right-hand side of the condition.
	Right Cond
}

func (*BinCond) isNode() {}

func (*BinCond) isExpr() {}

func (*BinCond) isCond() {}
