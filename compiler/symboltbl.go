package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int

	Outer *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable{
	return &SymbolTable{
		store: make(map[string]Symbol),
		Outer: outer,
	}
}
func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: st.numDefinitions}
	if st.Outer == nil{
		symbol.Scope = GlobalScope
	}else {
		symbol.Scope = LocalScope
	}
	st.numDefinitions++
	st.store[name] = symbol
	return symbol
}
func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	if !ok && st.Outer != nil{
		sym,ok = st.Outer.Resolve(name)
		return sym,ok
	}
	return sym, ok
}
