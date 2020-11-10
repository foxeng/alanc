package semantic

import "fmt"

// Check performs the semantic checks on the provided AST.
func Check(ast *Ast) error {
	st := NewSymTab()
	if _, err := ast.Program.check(st); err != nil {
		return err
	}
	return nil
}

func (n *FuncDef) check(st *SymTab) (Type, error) {
	// TODO: Can main have parameters?
	// TODO: Can main have non-proc return?

	// Add ID to scope.
	if !st.Add(n) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}
	// Enter scope.
	st.Enter()
	// Descend on parameters.
	for _, p := range n.Parameters {
		if _, err := p.check(st); err != nil {
			return nil, err
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

	return nil, nil
}

func (n *ParDef) check(st *SymTab) (Type, error) {
	// Descend on VarDef.
	if _, err := n.VarDef.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *PrimVarDef) check(st *SymTab) (Type, error) {
	// Add ID to scope.
	if !st.Add(n) {
		return nil, fmt.Errorf("%q already defined", n.ID)
	}

	return nil, nil
}

func (n *ArrayDef) check(st *SymTab) (Type, error) {
	// Descend on PrimVarDef.
	if _, err := n.PrimVarDef.check(st); err != nil {
		return nil, err
	}

	return nil, nil
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
	if _, err := n.Left.check(st); err != nil {
		return nil, err
	}
	// Descend on r-value.
	if _, err := n.Right.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *FuncCall) check(st *SymTab) (Type, error) {
	// Lookup ID.
	if st.Lookup(n.ID) == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}
	// Descend on args.
	for _, a := range n.Args {
		if _, err := a.check(st); err != nil {
			return nil, err
		}
	}

	return nil, nil
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
	if _, err := n.Expr.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *IntConstExpr) check(st *SymTab) (Type, error) {
	// Nothing to do
	return nil, nil
}

func (n *CharConstExpr) check(st *SymTab) (Type, error) {
	// Nothing to do
	return nil, nil
}

func (n *VarRef) check(st *SymTab) (Type, error) {
	// Lookup ID.
	if st.Lookup(n.ID) == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}

	return nil, nil
}

func (n *ArrayElem) check(st *SymTab) (Type, error) {
	// Lookup ID.
	if st.Lookup(n.ID) == nil {
		return nil, fmt.Errorf("%q not defined", n.ID)
	}
	// Descend on expression.
	if _, err := n.Index.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *StrLitExpr) check(st *SymTab) (Type, error) {
	// Nothing to do
	return nil, nil
}

func (n *FuncCallExpr) check(st *SymTab) (Type, error) {
	// Descend on FuncCall.
	if _, err := n.FuncCall.check(st); err != nil {
		return nil, err
	}
	return nil, nil
}

func (n *UnArithExpr) check(st *SymTab) (Type, error) {
	// Descend on expression.
	if _, err := n.Expr.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *BinArithExpr) check(st *SymTab) (Type, error) {
	// Descend on Left.
	if _, err := n.Left.check(st); err != nil {
		return nil, err
	}
	// Descend on Right.
	if _, err := n.Right.check(st); err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *ConstCond) check(st *SymTab) (Type, error) {
	// Nothing to do
	return nil, nil
}

func (n *UnCond) check(st *SymTab) (Type, error) {
	// Descend on condition.
	if _, err := n.Cond.check(st); err != nil {
		return nil, err
	}

	return nil, nil
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

	return nil, nil
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

	return nil, nil
}
