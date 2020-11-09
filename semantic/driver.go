package semantic

import "fmt"

// Check performs the semantic checks on the provided AST.
func Check(ast *Ast) error {
	st := NewSymTab()
	return symCheck(st, ast.Program)
}

// symCheck performs the elementary symbol checks on the subtree rooted in node.
func symCheck(st *SymTab, node Node) error {
	switch n := node.(type) {
	case *FuncDef:
		// Add ID to scope.
		if !st.Add(n) {
			return fmt.Errorf("%q already defined", n.ID)
		}
		// Enter scope.
		st.Enter()
		// Descend on parameters.
		for _, p := range n.Parameters {
			if err := symCheck(st, &p); err != nil {
				return err
			}
		}
		// Descend on local definitions.
		for _, ld := range n.LDefs {
			if err := symCheck(st, ld); err != nil {
				return err
			}
		}
		// Descend on body.
		if err := symCheck(st, &n.CompStmt); err != nil {
			return err
		}
		// Exit scope.
		st.Exit()
	case *ParDef:
		// Descend on VarDef.
		if err := symCheck(st, n.VarDef); err != nil {
			return err
		}
	case *PrimVarDef:
		// Add ID to scope.
		if !st.Add(n) {
			return fmt.Errorf("%q already defined", n.ID)
		}
	case *ArrayDef:
		// Descend on PrimVarDef.
		if err := symCheck(st, &n.PrimVarDef); err != nil {
			return err
		}
	case *CompStmt:
		// Descend on each statement.
		for _, s := range n.Stmts {
			if err := symCheck(st, s); err != nil {
				return err
			}
		}
	case *AssignStmt:
		// Descend on l-value.
		if err := symCheck(st, n.Left); err != nil {
			return err
		}
		// Descend on r-value.
		if err := symCheck(st, n.Right); err != nil {
			return err
		}
	case *FuncCall:
		// Lookup ID.
		if st.Lookup(n.ID) == nil {
			return fmt.Errorf("%q not defined", n.ID)
		}
		// Descend on args.
		for _, a := range n.Args {
			if err := symCheck(st, a); err != nil {
				return err
			}
		}
	case *FuncCallStmt:
		// Descend on FuncCall.
		if err := symCheck(st, &n.FuncCall); err != nil {
			return err
		}
	case *IfStmt:
		// Descend on condition.
		if err := symCheck(st, n.Cond); err != nil {
			return err
		}
		// Descend on statement.
		if err := symCheck(st, n.Stmt); err != nil {
			return err
		}
	case *IfElseStmt:
		// Descend on condition.
		if err := symCheck(st, n.Cond); err != nil {
			return err
		}
		// Descend on if branch statement.
		if err := symCheck(st, n.Stmt1); err != nil {
			return err
		}
		// Descend on else branch statement.
		if err := symCheck(st, n.Stmt2); err != nil {
			return err
		}
	case *WhileStmt:
		// Descend on condition.
		if err := symCheck(st, n.Cond); err != nil {
			return err
		}
		// Descend on statement.
		if err := symCheck(st, n.Stmt); err != nil {
			return err
		}
	case *ReturnStmt:
		if n.Expr == nil {
			return nil
		}
		// Descend on expression.
		if err := symCheck(st, n.Expr); err != nil {
			return err
		}
	case *IntConstExpr:
		return nil
	case *CharConstExpr:
		return nil
	case *VarRef:
		// Lookup ID.
		if st.Lookup(n.ID) == nil {
			return fmt.Errorf("%q not defined", n.ID)
		}
	case *ArrayElem:
		// Lookup ID.
		if st.Lookup(n.ID) == nil {
			return fmt.Errorf("%q not defined", n.ID)
		}
		// Descend on expression.
		if err := symCheck(st, n.Index); err != nil {
			return err
		}
	case *StrLitExpr:
		return nil
	case *FuncCallExpr:
		// Descend on FuncCall.
		if err := symCheck(st, &n.FuncCall); err != nil {
			return err
		}
	case *UnArithExpr:
		// Descend on expression.
		if err := symCheck(st, n.Expr); err != nil {
			return err
		}
	case *BinArithExpr:
		// Descend on Left.
		if err := symCheck(st, n.Left); err != nil {
			return err
		}
		// Descend on Right.
		if err := symCheck(st, n.Right); err != nil {
			return err
		}
	case *ConstCond:
		return nil
	case *UnCond:
		// Descend on condition.
		if err := symCheck(st, n.Cond); err != nil {
			return err
		}
	case *CompCond:
		// Descend on Left.
		if err := symCheck(st, n.Left); err != nil {
			return err
		}
		// Descend on Right.
		if err := symCheck(st, n.Right); err != nil {
			return err
		}
	case *BinCond:
		// Descend on Left.
		if err := symCheck(st, n.Left); err != nil {
			return err
		}
		// Descend on Right.
		if err := symCheck(st, n.Right); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("invalid node type: %T", node))
	}
	return nil
}
