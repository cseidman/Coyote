package symbol

import (
	. "../common"
)

type Symbol struct {
	Name       string
	ScopeDepth int
	Index      int
	Type       VarType
	Id         int
}

type SymbolTable struct {
	SymbolIndex int
	Symbols     []Symbol
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		SymbolIndex: 0,
		Symbols:     make([]Symbol, 65000),
	}
}

func (s *SymbolTable) AddSymbol(name string, scopedepth int, varType VarType, ContextId int) int {
	symb := Symbol{
		Name:       name,
		ScopeDepth: scopedepth,
		Index:      s.SymbolIndex,
		Type:       varType,
		Id:         ContextId,
	}
	s.Symbols[s.SymbolIndex] = symb
	s.SymbolIndex++
	return s.SymbolIndex - 1
}

func (s *SymbolTable) FindSymbol(name string, scopedepth int, varType VarType, contextId int) int {
	// Find the value in the symboltable
	for i := 0; i < s.SymbolIndex; i++ {
		if s.Symbols[i].Name == name &&
			s.Symbols[i].ScopeDepth == scopedepth &&
			s.Symbols[i].Type == varType &&
			s.Symbols[i].Id == contextId {
			return i
		}
	}
	return -1
}
