package assembler

import (
	"fmt"
)

type SymbolTable struct {
	addressMap  map[string]int
	nextAddress int
}

func NewSymbolTable() *SymbolTable {
	symbolTable := &SymbolTable{
		addressMap: map[string]int{
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SCREEN": 16384,
			"KBD":    24576,
		},
		nextAddress: 16,
	}
	return symbolTable
}

func (st *SymbolTable) AddEntry(symbol string, address int) {
	st.addressMap[symbol] = address
}

func (st SymbolTable) Contains(symbol string) bool {
	if _, ok := st.addressMap[symbol]; ok {
		return true
	}
	return false
}

func (st SymbolTable) GetAddress(symbol string) int {
	return st.addressMap[symbol]
}

func (st SymbolTable) GetNextAddress() int {
	return st.nextAddress
}

func (st *SymbolTable) IncrementNextAddress() {
	st.nextAddress++
}

func (st SymbolTable) OutputSymbolTable() {
	for k, v := range st.addressMap {
		fmt.Println(k, v)
	}
}
