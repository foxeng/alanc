package semantic

import "testing"

func TestNewSymTab(t *testing.T) {
	st := NewSymTab()
	if st == nil {
		t.Error("NewSymTab() = nil")
	}
	st.Enter()
	// All standard library functions should be predefined.
	for _, f := range stdlib {
		if st.Lookup(f.Id()) == nil {
			t.Errorf("%q not found in fresh main scope", f.Id())
		}
	}
	// TODO OPT: Test that nothing but the stdlib definitions is predefined.
}

func TestSymTabEnter(t *testing.T) {
	st := NewSymTab()
	st.Enter()

	// Enter() should increment the stack depth.
	ol := len(st.scopes.stack)
	st.Enter()
	nl := len(st.scopes.stack)
	if nl != ol+1 {
		t.Errorf("Enter() changed scope stack length by %d, want 1", nl-ol)
	}
}

func TestSymTab(t *testing.T) {
	st := NewSymTab()
	st.Enter()

	defs := []LocalDef{
		&FuncDef{
			ID: "testFunc",
		},
		&ParDef{
			VarDef: &PrimVarDef{
				ID:       "testPar",
				DataType: DataTypeInt,
			},
		},
		&PrimVarDef{
			ID:       "testPrimVar",
			DataType: DataTypeByte,
		},
		&ArrayDef{
			PrimVarDef: PrimVarDef{
				ID:       "testArray",
				DataType: DataTypeInt,
			},
		},
	}
	for _, d := range defs {
		r := st.Lookup(d.Id())
		if r != nil {
			t.Errorf("%q found in scope before its addition", d.Id())
		}
		if !st.Add(d) {
			t.Errorf("failed to add %q", d.Id())
		}
		r = st.Lookup(d.Id())
		if r == nil {
			t.Errorf("%q not found in scope after its addition", d.Id())
		}
	}

	// Standard library symbols should be shadowable.
	stdShadowDef := &PrimVarDef{
		ID:       "extend",
		DataType: DataTypeInt,
	}
	r := st.Lookup(stdShadowDef.Id())
	if _, ok := r.(*FuncDef); !ok {
		t.Errorf("predefined %q not a FuncDef", stdShadowDef.Id())
	}
	if !st.Add(stdShadowDef) {
		t.Errorf("failed to add %q", stdShadowDef.Id())
	}
	r = st.Lookup(stdShadowDef.Id())
	if _, ok := r.(*PrimVarDef); !ok {
		t.Errorf("shadowed %q not a PrimVarDef", stdShadowDef.Id())
	}

	// Add() and Lookup() should work when in a nested scope.
	st.Enter()
	r = st.Lookup(defs[0].Id())
	if r == nil {
		t.Errorf("%q not found in inner scope", defs[0].Id())
	}
	def2 := &FuncDef{
		ID: "testFunc2",
	}
	r = st.Lookup(def2.Id())
	if r != nil {
		t.Errorf("%q found in scope before its addition", def2.Id())
	}
	if !st.Add(def2) {
		t.Errorf("failed to add %q", def2.Id())
	}
	st.Enter()
	r = st.Lookup(def2.Id())
	if r == nil {
		t.Errorf("%q not found in scope after its addition", def2.Id())
	}

	// Symbols from outer scopes should be shadowable.
	shadowDef := &PrimVarDef{
		ID:       "testPrimVar",
		DataType: DataTypeInt,
	}
	r = st.Lookup(shadowDef.Id())
	odt := r.(*PrimVarDef).DataType
	if !st.Add(shadowDef) {
		t.Errorf("failed to add %q", shadowDef.Id())
	}
	r = st.Lookup(shadowDef.Id())
	ndt := r.(*PrimVarDef).DataType
	if ndt == odt {
		t.Errorf("%q not shadowed", shadowDef.Id())
	}

	// Symbols from current scope should not be shadowable (no redefinitions).
	if st.Add(shadowDef) {
		t.Errorf("%q redefined in the same scope", shadowDef.Id())
	}
}

func TestSymTabExit(t *testing.T) {
	st := NewSymTab()
	st.Enter()

	def := &FuncDef{
		ID: "testFunc",
	}
	st.Enter()
	st.Add(def)

	// Exit() should decrement the stack depth.
	ol := len(st.scopes.stack)
	st.Exit()
	nl := len(st.scopes.stack)
	if nl != ol-1 {
		t.Errorf("Exit() changed scope stack length by %d, want -1", nl-ol)
	}
	// TODO OPT: Test corner case (empty stack).

	// Exit() should remove all symbols defined in the current scope.
	r := st.Lookup(def.Id())
	if r != nil {
		t.Errorf("%q found in scope after its removal", def.ID)
	}
}
