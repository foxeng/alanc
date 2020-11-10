package semantic

// Pascal scope: each symbol is visible from the point of its declaration until the end of that
// unit. Except if it's shadowed.

// scope is a single Alan scope (the scope of a unit, not a single symbol).
type scope map[ID]Type

// scopeStack is a stack of scopes.
type scopeStack struct {
	stack []scope
}

// push pushes a new empty scope to the top of the stack.
func (st *scopeStack) push() {
	st.stack = append(st.stack, scope{})
}

// pop pops the scope at the top of the stack.
func (st *scopeStack) pop() {
	l := len(st.stack)
	if l == 1 {
		panic("popped the standard library (pre-main) scope")
	}
	// TODO OPT: Optimize for space, this will never free any underlying memory.
	st.stack = st.stack[:l-1]
}

// top returns the current top of the stack.
func (st *scopeStack) top() *scope {
	return &st.stack[len(st.stack)-1]
}

// SymTab is the Symbol Table for Alan.
type SymTab struct {
	// TODO OPT: Make lookup constant: use a single map for the current scope. Update it when
	// entering / exiting scopes and defining new symbols.

	// scopes is the stack of scopes. This provides linear lookup, but constant addition, enter and
	// exit, but most importantly a simple implementation.
	scopes scopeStack
}

// NewSymTab returns a new Symbol Table. That is left in the standard library (pre-main) scope, so
// Enter() should be called on it before any further use, to enter the main program scope.
func NewSymTab() *SymTab {
	st := &SymTab{}
	st.Enter()
	// Inject standard library definitions in the outermost scope (nothing else should be defined
	// in that, so as for them to be immediately shadowable, from the outermost program scope).
	for _, fd := range stdlib {
		st.Add(fd.ID, fd.FunctionType)
	}
	return st
}

// Enter creates a new scope and switches to it.
func (st *SymTab) Enter() {
	st.scopes.push()
}

// Add adds a new symbol definition for name to the current scope, returning false if there is a
// definition for that name in the current scope already (not shadowable).
func (st *SymTab) Add(name ID, t Type) bool {
	sc := st.scopes.top()
	if _, ok := (*sc)[name]; ok {
		return false
	}

	(*sc)[name] = t
	return true
}

// Lookup searches if name is visible from the current scope. If so, it returns its type, otherwise
// it returns nil.
func (st *SymTab) Lookup(name ID) Type {
	for i := len(st.scopes.stack) - 1; i >= 0; i-- {
		if t, ok := st.scopes.stack[i][name]; ok {
			return t
		}
	}
	return nil
}

// Exit removes the current scope and switches to its previous.
func (st *SymTab) Exit() {
	st.scopes.pop()
}
