package semantic

import "testing"

func TestNewSymTab(t *testing.T) {
	st := NewSymTab()
	if st == nil {
		t.Error("NewSymTab() = nil")
	}
	st.Enter("")
	// All standard library functions should be predefined.
	for _, f := range stdlib {
		if st.Lookup(f.ID) == nil {
			t.Errorf("%q not found in fresh main scope", f.ID)
		}
	}
	// TODO OPT: Test that nothing but the stdlib definitions is predefined.
}

func TestSymTabEnter(t *testing.T) {
	st := NewSymTab()
	st.Enter("")

	// Enter() should increment the stack depth.
	ol := len(st.scopes.stack)
	st.Enter("")
	nl := len(st.scopes.stack)
	if nl != ol+1 {
		t.Errorf("Enter() changed scope stack length by %d, want 1", nl-ol)
	}
}

func TestSymTabCurrentID(t *testing.T) {
	st := NewSymTab()
	testID := ID("testID")
	st.Enter(testID)

	cid := st.CurrentID()
	if cid != testID {
		t.Errorf("CurrentID() returned %q, want %q", cid, testID)
	}
}

func TestSymTab(t *testing.T) {
	st := NewSymTab()
	st.Enter("")

	defs := []struct {
		ID
		Type
	}{
		{
			ID: "testFunc",
			Type: FunctionType{
				Parameters: []ParameterType{},
			},
		},
		{
			ID: "testPar",
			Type: ParameterType{
				DType: PrimitiveTypeInt,
			},
		},
		{
			ID:   "testPrimVar",
			Type: PrimitiveTypeByte,
		},
		{
			ID: "testArray",
			Type: ArrayType{
				PrimitiveType: PrimitiveTypeInt,
			},
		},
	}
	for _, d := range defs {
		r := st.Lookup(d.ID)
		if r != nil {
			t.Errorf("%q found in scope before its addition", d.ID)
		}
		if !st.Add(d.ID, d.Type) {
			t.Errorf("failed to add %q", d.ID)
		}
		r = st.Lookup(d.ID)
		if r == nil {
			t.Errorf("%q not found in scope after its addition", d.ID)
		}
	}

	// Standard library symbols should be shadowable.
	stdShadowDef := struct {
		ID
		PrimitiveType
	}{
		ID:            "extend",
		PrimitiveType: PrimitiveTypeInt,
	}
	r := st.Lookup(stdShadowDef.ID)
	if _, ok := r.(FunctionType); !ok {
		t.Errorf("predefined %q not a FunctionType", stdShadowDef.ID)
	}
	if !st.Add(stdShadowDef.ID, stdShadowDef.PrimitiveType) {
		t.Errorf("failed to add %q", stdShadowDef.ID)
	}
	r = st.Lookup(stdShadowDef.ID)
	if _, ok := r.(PrimitiveType); !ok {
		t.Errorf("shadowed %q not a PrimitiveType", stdShadowDef.ID)
	}

	// Add() and Lookup() should work when in a nested scope.
	st.Enter("")
	r = st.Lookup(defs[0].ID)
	if r == nil {
		t.Errorf("%q not found in inner scope", defs[0].ID)
	}
	def2 := struct {
		ID
		FunctionType
	}{
		ID: "testFunc2",
		FunctionType: FunctionType{
			Parameters: []ParameterType{},
		},
	}
	r = st.Lookup(def2.ID)
	if r != nil {
		t.Errorf("%q found in scope before its addition", def2.ID)
	}
	if !st.Add(def2.ID, def2.FunctionType) {
		t.Errorf("failed to add %q", def2.ID)
	}
	st.Enter("")
	r = st.Lookup(def2.ID)
	if r == nil {
		t.Errorf("%q not found in scope after its addition", def2.ID)
	}

	// Symbols from outer scopes should be shadowable.
	shadowDef := struct {
		ID
		PrimitiveType
	}{
		ID:            "testPrimVar",
		PrimitiveType: PrimitiveTypeInt,
	}
	r = st.Lookup(shadowDef.ID)
	odt := r.(PrimitiveType)
	if !st.Add(shadowDef.ID, shadowDef.PrimitiveType) {
		t.Errorf("failed to add %q", shadowDef.ID)
	}
	r = st.Lookup(shadowDef.ID)
	ndt := r.(PrimitiveType)
	if ndt == odt {
		t.Errorf("%q not shadowed", shadowDef.ID)
	}

	// Symbols from current scope should not be shadowable (no redefinitions).
	if st.Add(shadowDef.ID, shadowDef.PrimitiveType) {
		t.Errorf("%q redefined in the same scope", shadowDef.ID)
	}
}

func TestSymTabExit(t *testing.T) {
	st := NewSymTab()
	st.Enter("")

	def := struct {
		ID
		FunctionType
	}{
		ID: "testFunc",
		FunctionType: FunctionType{
			Parameters: []ParameterType{},
		},
	}
	st.Enter("")
	st.Add(def.ID, def.FunctionType)

	// Exit() should decrement the stack depth.
	ol := len(st.scopes.stack)
	st.Exit()
	nl := len(st.scopes.stack)
	if nl != ol-1 {
		t.Errorf("Exit() changed scope stack length by %d, want -1", nl-ol)
	}
	// TODO OPT: Test corner case (empty stack).

	// Exit() should remove all symbols defined in the current scope.
	r := st.Lookup(def.ID)
	if r != nil {
		t.Errorf("%q found in scope after its removal", def.ID)
	}
}
