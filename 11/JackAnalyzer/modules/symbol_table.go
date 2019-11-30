package modules

import (
	// "fmt"
	"strconv"
)

type SymbolTable struct {
	classSymbolTable           map[string]map[string]string
	classSymbolTableIndex      int
	subroutineSymbolTable      map[string]map[string]string
	subroutineSymbolTableIndex int
}

func (self *SymbolTable) initialize() {
	self.classSymbolTable = make(map[string]map[string]string)
	self.classSymbolTableIndex = 0
	self.subroutineSymbolTable = make(map[string]map[string]string)
	self.subroutineSymbolTableIndex = 0
}

func NewSymbolTable() *SymbolTable {
	symbolTable := &SymbolTable{}
	symbolTable.initialize()
	return symbolTable
}

/**
 * 新しいサブルーチンスコープを開始（サブルーチンのたびに呼び出せば良い）
 */
func (self *SymbolTable) StartSubroutine() {
	self.subroutineSymbolTable = make(map[string]map[string]string)
	self.subroutineSymbolTableIndex = 0
}

/**
 * 新しい識別子を定義・実行インデックスの割り当て
 */
func (self *SymbolTable) Define(name string, typeName string, kind string) {
	if self.isClassScope(kind) {
		self.classSymbolTable[name] = make(map[string]string)
		self.classSymbolTable[name]["typeName"] = typeName
		self.classSymbolTable[name]["kind"] = kind
		self.classSymbolTable[name]["index"] = strconv.Itoa(self.classSymbolTableIndex)
		self.classSymbolTableIndex++
	} else if self.isSubroutineScope(kind) {
		self.subroutineSymbolTable[name] = make(map[string]string)
		self.subroutineSymbolTable[name]["typeName"] = typeName
		self.subroutineSymbolTable[name]["kind"] = kind
		self.subroutineSymbolTable[name]["index"] = strconv.Itoa(self.subroutineSymbolTableIndex)
		self.subroutineSymbolTableIndex++
	}
}

func (self *SymbolTable) isClassScope(kind string) bool {
	return kind == "static" || kind == "field"
}

func (self *SymbolTable) isSubroutineScope(kind string) bool {
	return kind == "arg" || kind == "var"
}

func (self *SymbolTable) VarCount(kind string) int {
	result := 0

	symbolTable := self.classSymbolTable
	if self.isSubroutineScope(kind) {
		symbolTable = self.subroutineSymbolTable
	}

	for _, symbol := range symbolTable {
		if symbol["kind"] == kind {
			result++
		}
	}

	return result
}

func (self *SymbolTable) KindOf(name string) string {
	if symbol, ok := self.classSymbolTable[name]; ok {
		return symbol["kind"]
	} else {
		return self.subroutineSymbolTable[name]["kind"]
	}
}

func (self *SymbolTable) TypeOf(name string) string {
	if symbol, ok := self.classSymbolTable[name]; ok {
		return symbol["typeName"]
	} else {
		return self.subroutineSymbolTable[name]["typeName"]
	}
}

func (self *SymbolTable) IndexOf(name string) int {
	if symbol, ok := self.classSymbolTable[name]; ok {
		index, _ := strconv.Atoi(symbol["index"])
		return index
	} else {
		index, _ := strconv.Atoi(self.subroutineSymbolTable[name]["index"])
		return index
	}
}
