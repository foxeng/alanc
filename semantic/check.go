package semantic

import (
	"errors"
	"fmt"
)

// Check performs the semantic checks on the provided AST.
func Check(ast *Ast) error {
	st := NewSymTab()
	if _, err := ast.Program.check(st); err != nil {
		return err
	}
	return nil
}

func (n *FuncDef) check(st *SymTab) (Type, error) {
	// NOTE: Ideally, we would add the function to the current scope, enter a new scope and proceed
	// with the rest (parameters, locals, etc.). But, to add the function we need to know the
	// parameter and return types. To do this, we enter a new, temporary scope just to get the
	// parameters, then exit it, and proceed as normal.
	fType := FunctionType{
		Parameters: make([]ParameterType, len(n.Parameters)),
		Return:     n.RType,
	}
	// Enter temporary scope.
	st.Enter("")
	// Descend on parameters.
	for i, p := range n.Parameters {
		t, err := p.check(st)
		if err != nil {
			return nil, err
		}
		fType.Parameters[i] = t.(ParameterType)
	}
	// Exit temporary scope.
	st.Exit()

	// Check that main has no parameters and has proc return type.
	if st.CurrentID() == "" {
		if len(fType.Parameters) > 0 {
			return nil, errors.New("main function cannot accept parameters")
		}
		if fType.Return != nil {
			return nil, errors.New("main function must have proc return type")
		}
	}

	// Add to scope.
	if !st.Add(n.ID, fType) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}
	// Enter scope.
	st.Enter(n.ID)
	// Descend on parameters (just to add them to the scope).
	for _, p := range n.Parameters {
		_, err := p.check(st)
		if err != nil {
			// NOTE: This should never happen, they have been checked above.
			panic(fmt.Sprintf("error not already caught: %v", err))
		}
	}
	// Descend on local definitions.
	for _, ld := range n.LDefs {
		if _, err := ld.check(st); err != nil {
			return nil, err
		}
	}
	// Descend on body.
	if _, err := n.CompStmt.check(st); err != nil {
		return nil, err
	}
	// Exit scope.
	st.Exit()

	return fType, nil
}

func (n *ParDef) check(st *SymTab) (Type, error) {
	// Add to scope.
	if !st.Add(n.ID, n.Type.DType) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}

	return n.Type, nil
}

func (n *PrimVarDef) check(st *SymTab) (Type, error) {
	// Add to scope.
	if !st.Add(n.ID, n.Type) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}

	return n.Type, nil
}

func (n *ArrayDef) check(st *SymTab) (Type, error) {
	// Add to scope.
	if !st.Add(n.ID, n.Type) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}

	return n.Type, nil
}

func (n *CompStmt) check(st *SymTab) (Type, error) {
	// Descend on each statement.
	for _, s := range n.Stmts {
		if _, err := s.check(st); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (n *AssignStmt) check(st *SymTab) (Type, error) {
	// Descend on l-value.
	lt, err := n.Left.check(st)
	if err != nil {
		return nil, err
	}
	// Descend on r-value.
	rt, err := n.Right.check(st)
	if err != nil {
		return nil, err
	}

	// Check l-value and r-value are of the same, primitive type.
	plt, ok := lt.(PrimitiveType)
	if !ok {
		return nil, fmt.Errorf("cannot assign to non-primitive type %T", lt)
	}
	prt, ok := rt.(PrimitiveType)
	if !ok {
		return nil, fmt.Errorf("cannot assign non-primitive type %T", lt)
	}
	if plt != prt {
		return nil, fmt.Errorf("cannot assign primitive type %T to %T", prt, plt)
	}

	return nil, nil
}

func (n *FuncCall) check(st *SymTab) (Type, error) {
	// Lookup ID.
	t := st.Lookup(n.ID)
	if t == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}
	ft, ok := t.(FunctionType)
	if !ok {
		return nil, fmt.Errorf("%q not a function", n.ID)
	}

	// Descend on arguments.
	for i, a := range n.Args {
		t, err := a.check(st)
		if err != nil {
			return nil, err
		}
		// Check argument type-matches corresponding parameter.
		switch pt := ft.Parameters[i].DType.(type) {
		case PrimitiveType:
			if t != pt {
				return nil, fmt.Errorf("argument #%d to %q has type %T, want %T", i+1, n.ID, t, pt)
			}
		case ArrayType:
			// Ignore array sizes (i.e. only check the element types match).
			at, ok := t.(ArrayType)
			if !ok || at.PrimitiveType != pt.PrimitiveType {
				return nil, fmt.Errorf("argument #%d to %q has type %T, want %T", i+1, n.ID, t, pt)
			}
		default:
			panic(fmt.Sprintf("function parameter of invalid data type %T", pt))
		}
		if ft.Parameters[i].IsRef {
			// Check argument is an l-value.
			if _, ok := a.(LVal); !ok {
				return nil, fmt.Errorf("argument #%d to %q cannot be passed by reference (not an "+
					"l-value", i+1, n.ID)
			}
		}
	}

	return ft.Return, nil
}

func (n *FuncCallStmt) check(st *SymTab) (Type, error) {
	// Descend on FuncCall.
	if _, err := n.FuncCall.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *IfStmt) check(st *SymTab) (Type, error) {
	// Descend on condition.
	// TODO OPT: Check condition type? Shouldn't be necessary...
	if _, err := n.Cond.check(st); err != nil {
		return nil, err
	}
	// Descend on statement.
	if _, err := n.Stmt.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *IfElseStmt) check(st *SymTab) (Type, error) {
	// Descend on condition.
	// TODO OPT: Check condition type? Shouldn't be necessary...
	if _, err := n.Cond.check(st); err != nil {
		return nil, err
	}
	// Descend on if branch statement.
	if _, err := n.Stmt1.check(st); err != nil {
		return nil, err
	}
	// Descend on else branch statement.
	if _, err := n.Stmt2.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *WhileStmt) check(st *SymTab) (Type, error) {
	// Descend on condition.
	// TODO OPT: Check condition type? Shouldn't be necessary...
	if _, err := n.Cond.check(st); err != nil {
		return nil, err
	}
	// Descend on statement.
	if _, err := n.Stmt.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *ReturnStmt) check(st *SymTab) (Type, error) {
	if n.Expr == nil {
		return nil, nil
	}
	// Descend on expression.
	t, err := n.Expr.check(st)
	if err != nil {
		return nil, err
	}
	// Check expression type matches return type of enclosing function.
	fName := st.CurrentID()
	fRet := st.Lookup(fName).(FunctionType).Return
	if fRet == nil {
		if t != nil {
			return nil, fmt.Errorf("return %T from procedure %q", t, fName)
		}
	}
	if t != *fRet {
		return nil, fmt.Errorf("return %T from function %q with return type %T", t, fName, fRet)
	}

	return nil, nil
}

func (n *IntConstExpr) check(st *SymTab) (Type, error) {
	return PrimitiveTypeInt, nil
}

func (n *CharConstExpr) check(st *SymTab) (Type, error) {
	return PrimitiveTypeByte, nil
}

func (n *VarRef) check(st *SymTab) (Type, error) {
	// Lookup ID.
	t := st.Lookup(n.ID)
	if t == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}
	pt, ok := t.(DType)
	if !ok {
		return nil, fmt.Errorf("%q not a variable", n.ID)
	}

	return pt, nil
}

func (n *ArrayElem) check(st *SymTab) (Type, error) {
	// Lookup ID.
	t := st.Lookup(n.ID)
	if t == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}
	at, ok := t.(ArrayType)
	if !ok {
		return nil, fmt.Errorf("%q not an array", n.ID)
	}

	// Descend on expression.
	t, err := n.Index.check(st)
	if err != nil {
		return nil, err
	}
	// Check expression is int.
	switch et := t.(type) {
	case PrimitiveType:
		if et != PrimitiveTypeInt {
			return nil, fmt.Errorf("array index of primitive type %T, need \"int\"", et)
		}
	default:
		return nil, fmt.Errorf("array index of non-primitive type %T, need \"int\"", t)
	}

	return at.PrimitiveType, nil
}

func (n *StrLitExpr) check(st *SymTab) (Type, error) {
	return ArrayType{
		PrimitiveType: PrimitiveTypeByte,
		Size:          len(n.Val) + 1,
	}, nil
}

func (n *FuncCallExpr) check(st *SymTab) (Type, error) {
	// Descend on FuncCall.
	t, err := n.FuncCall.check(st)
	if err != nil {
		return nil, err
	}
	rt := t.(*PrimitiveType)
	if rt == nil {
		return nil, fmt.Errorf("cannot call procedure %q in an expression", n.ID)
	}

	return *rt, nil
}

func (n *UnArithExpr) check(st *SymTab) (Type, error) {
	// Descend on expression.
	t, err := n.Expr.check(st)
	if err != nil {
		return nil, err
	}
	// Check expression is int.
	switch et := t.(type) {
	case PrimitiveType:
		if et != PrimitiveTypeInt {
			return nil, fmt.Errorf("unary arithmetic expression of primitive type %T, need \"int\"",
				et)
		}
	default:
		return nil, fmt.Errorf("unary arithmetic expression of non-primitive type %T, need \"int\"",
			t)
	}

	return PrimitiveTypeInt, nil
}

func (n *BinArithExpr) check(st *SymTab) (Type, error) {
	// Descend on Left.
	lt, err := n.Left.check(st)
	if err != nil {
		return nil, err
	}
	// Descend on Right.
	rt, err := n.Right.check(st)
	if err != nil {
		return nil, err
	}

	// Check Left and Right type-match (int or byte).
	plt, ok := lt.(PrimitiveType)
	if !ok {
		return nil, fmt.Errorf("left operand of binary arithmetic expression of non-primitive "+
			"type %T", lt)
	}
	prt, ok := rt.(PrimitiveType)
	if !ok {
		return nil, fmt.Errorf("right operand of binary arithmetic expression of non-primitive "+
			"type %T", lt)
	}
	if plt != prt {
		return nil, fmt.Errorf("cannot apply binary arithmetic operator to primitive types %T and "+
			"%T", plt, prt)
	}

	return plt, nil
}

func (n *ConstCond) check(st *SymTab) (Type, error) {
	return PrimitiveTypeBool, nil
}

func (n *UnCond) check(st *SymTab) (Type, error) {
	// Descend on condition.
	if _, err := n.Cond.check(st); err != nil {
		return nil, err
	}

	return PrimitiveTypeBool, nil
}

func (n *CompCond) check(st *SymTab) (Type, error) {
	// Descend on Left.
	if _, err := n.Left.check(st); err != nil {
		return nil, err
	}
	// Descend on Right.
	if _, err := n.Right.check(st); err != nil {
		return nil, err
	}

	return PrimitiveTypeBool, nil
}

func (n *BinCond) check(st *SymTab) (Type, error) {
	// Descend on Left.
	if _, err := n.Left.check(st); err != nil {
		return nil, err
	}
	// Descend on Right.
	if _, err := n.Right.check(st); err != nil {
		return nil, err
	}

	return PrimitiveTypeBool, nil
}
