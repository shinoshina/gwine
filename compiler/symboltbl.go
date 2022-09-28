package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
	Outer       *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:       make(map[string]Symbol),
		FreeSymbols: []Symbol{},
	}
}
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		store:       make(map[string]Symbol),
		FreeSymbols: []Symbol{},
		Outer:       outer,
	}
}
func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: st.numDefinitions}
	if st.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	st.numDefinitions++
	st.store[name] = symbol
	return symbol
}
func (st *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	st.store[name] = symbol
	return symbol
}
func (st *SymbolTable) defineFree(original Symbol) Symbol {
	st.FreeSymbols = append(st.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(st.FreeSymbols) - 1, Scope: FreeScope}
	st.store[original.Name] = symbol
	return symbol
}
func (st *SymbolTable) DefineFunctionName(name string) Symbol{
	symbol := Symbol{Name: name,Index: 0,Scope: FunctionScope}
	st.store[name] = symbol
	return symbol
}
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	if !ok && st.Outer != nil {
		sym, ok = st.Outer.Resolve(name)
		if !ok {
			return sym, ok
		}
		if sym.Scope == GlobalScope || sym.Scope == BuiltinScope {
			return sym, ok
		}

		return st.defineFree(sym), true
	}
	return sym, ok
}
