package modules

import (
	"fmt"
	"os"
	"strconv"
)

type CompilationEngine struct {
	tokenizer     *JackTokenizer
	symbolTable   *SymbolTable
	outputXmlFile *os.File
}

func (self *CompilationEngine) initialize(tokenizer *JackTokenizer, symbolTable *SymbolTable, outputXmlFile *os.File) {
	self.tokenizer = tokenizer
	self.symbolTable = symbolTable
	self.outputXmlFile = outputXmlFile
}

func NewCompilationEngine(tokenizer *JackTokenizer, symbolTable *SymbolTable, outputXmlFile *os.File) *CompilationEngine {
	compilationEngine := &CompilationEngine{}
	compilationEngine.initialize(tokenizer, symbolTable, outputXmlFile)
	return compilationEngine
}

func (self *CompilationEngine) CompileClass() {
	self.symbolTable.StartSubroutine()

	self.outputToXmlFile("<class>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // 'class'
	self.tokenizer.Advance()

	self.outputToXmlFile("<identifier category=class> " + self.tokenizer.Identifier() + " </identifier>") // className
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '{'
	self.tokenizer.Advance()

	for {
		if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "}" {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '}'
			break
		}

		if self.tokenizer.KeyWord() == "static" || self.tokenizer.KeyWord() == "field" {
			self.compileClassVarDec()
		} else if self.tokenizer.KeyWord() == "constructor" || self.tokenizer.KeyWord() == "function" || self.tokenizer.KeyWord() == "method" {
			self.compileSubroutineDec()
		}
	}

	self.outputToXmlFile("</class>")
}

func (self *CompilationEngine) compileClassVarDec() {
	self.outputToXmlFile("<classVarDec>")

	kind := self.tokenizer.KeyWord()
	self.outputToXmlFile("<keyword> " + kind + " </keyword>") // ('static'|'field')
	self.tokenizer.Advance()

	typeName := self.resolveType() // type
	self.tokenizer.Advance()

	name := self.tokenizer.Identifier()
	self.symbolTable.Define(name, typeName, kind)
	self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
	self.tokenizer.Advance()

	for {
		if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == ";" {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ';'
			self.tokenizer.Advance()
			break
		}

		if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "," {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ','
		} else if self.tokenizer.TokenType() == "IDENTIFIER" {
			name := self.tokenizer.Identifier()
			self.symbolTable.Define(name, typeName, kind)
			self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
		}
		self.tokenizer.Advance()
	}

	self.outputToXmlFile("</classVarDec>")
}

func (self *CompilationEngine) compileSubroutineDec() {
	self.symbolTable.StartSubroutine()

	self.outputToXmlFile("<subroutineDec>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // ('constructor'|'function'|'method')
	self.tokenizer.Advance()

	if self.tokenizer.TokenType() == "KEYWORD" && self.tokenizer.KeyWord() == "void" {
		self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // 'void'
	} else {
		self.resolveType() // type
	}
	self.tokenizer.Advance()

	self.outputToXmlFile("<identifier category=subroutine> " + self.tokenizer.Identifier() + " </identifier>") // subroutineName
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
	self.tokenizer.Advance()

	self.compileParameterList()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
	self.tokenizer.Advance()

	self.compileSubroutineBody()

	self.outputToXmlFile("</subroutineDec>")
}

func (self *CompilationEngine) compileParameterList() {
	self.outputToXmlFile("<parameterList>")
	if !self.isType() {
		self.outputToXmlFile("</parameterList>")
		return
	}

	typeName := self.resolveType() // type
	self.tokenizer.Advance()

	name := self.tokenizer.Identifier()
	self.symbolTable.Define(name, typeName, "arg")
	self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
	self.tokenizer.Advance()

	for {
		if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "," {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ','
			self.tokenizer.Advance()

			typeName := self.resolveType() // type
			self.tokenizer.Advance()

			name := self.tokenizer.Identifier()
			self.symbolTable.Define(name, typeName, "arg")
			self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
			self.tokenizer.Advance()
		} else {
			break
		}
	}

	self.outputToXmlFile("</parameterList>")
}

func (self *CompilationEngine) compileSubroutineBody() {
	self.outputToXmlFile("<subroutineBody>")

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '{'
	self.tokenizer.Advance()

	for {
		if self.tokenizer.TokenType() == "KEYWORD" && self.tokenizer.KeyWord() == "var" {
			self.compileVarDec()
		} else {
			break
		}
	}

	self.compileStatements()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '}'
	self.tokenizer.Advance()

	self.outputToXmlFile("</subroutineBody>")
}

func (self *CompilationEngine) compileVarDec() {
	self.outputToXmlFile("<varDec>")

	kind := self.tokenizer.KeyWord()
	self.outputToXmlFile("<keyword> " + kind + " </keyword>") // var
	self.tokenizer.Advance()

	typeName := self.resolveType() // type
	self.tokenizer.Advance()

	name := self.tokenizer.Identifier()
	self.symbolTable.Define(name, typeName, kind)
	self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
	self.tokenizer.Advance()

	for {
		if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "," {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ','
			self.tokenizer.Advance()

			name := self.tokenizer.Identifier()
			self.symbolTable.Define(name, typeName, kind)
			self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " defined> " + name + " </identifier>") // varName
			self.tokenizer.Advance()
		} else {
			break
		}
	}

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ';'
	self.tokenizer.Advance()

	self.outputToXmlFile("</varDec>")
}

func (self CompilationEngine) compileStatements() {
	self.outputToXmlFile("<statements>")

	for {
		if !self.isStatement() {
			break
		}
		if self.tokenizer.KeyWord() == "let" {
			self.compileLet()
		} else if self.tokenizer.KeyWord() == "if" {
			self.compileIf()
		} else if self.tokenizer.KeyWord() == "while" {
			self.compileWhile()
		} else if self.tokenizer.KeyWord() == "do" {
			self.compileDo()
		} else if self.tokenizer.KeyWord() == "return" {
			self.compileReturn()
		}
	}

	self.outputToXmlFile("</statements>")
}

func (self CompilationEngine) compileLet() {
	self.outputToXmlFile("<letStatement>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // let
	self.tokenizer.Advance()

	name := self.tokenizer.Identifier()
	self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " used> " + name + " </identifier>") // varName
	self.tokenizer.Advance()

	if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "[" {
		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '['
		self.tokenizer.Advance()

		self.compileExpression()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ']'
		self.tokenizer.Advance()
	}

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '='
	self.tokenizer.Advance()

	self.compileExpression()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ';'
	self.tokenizer.Advance()

	self.outputToXmlFile("</letStatement>")
}

func (self CompilationEngine) compileIf() {
	self.outputToXmlFile("<ifStatement>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // if
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
	self.tokenizer.Advance()

	self.compileExpression()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '{'
	self.tokenizer.Advance()

	self.compileStatements()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '}'
	self.tokenizer.Advance()

	if self.tokenizer.TokenType() == "KEYWORD" && self.tokenizer.KeyWord() == "else" {
		self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // else
		self.tokenizer.Advance()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '{'
		self.tokenizer.Advance()

		self.compileStatements()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '}'
		self.tokenizer.Advance()
	}

	self.outputToXmlFile("</ifStatement>")
}

func (self CompilationEngine) compileWhile() {
	self.outputToXmlFile("<whileStatement>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // while
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
	self.tokenizer.Advance()

	self.compileExpression()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '{'
	self.tokenizer.Advance()

	self.compileStatements()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '}'
	self.tokenizer.Advance()

	self.outputToXmlFile("</whileStatement>")
}

func (self CompilationEngine) compileDo() {
	self.outputToXmlFile("<doStatement>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // do
	self.tokenizer.Advance()

	self.resolveSubroutineCall(false)
	self.tokenizer.Advance()

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ';'
	self.tokenizer.Advance()

	self.outputToXmlFile("</doStatement>")
}

func (self CompilationEngine) compileReturn() {
	self.outputToXmlFile("<returnStatement>")

	self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>") // return
	self.tokenizer.Advance()

	if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == ";" {
		// 何もしない
	} else {
		self.compileExpression()
	}

	self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ';'
	self.tokenizer.Advance()

	self.outputToXmlFile("</returnStatement>")
}

func (self CompilationEngine) compileExpression() {
	self.outputToXmlFile("<expression>")

	self.compileTerm()

	for {
		if self.tokenizer.TokenType() == "SYMBOL" &&
			(self.tokenizer.Symbol() == "+" ||
				self.tokenizer.Symbol() == "-" ||
				self.tokenizer.Symbol() == "*" ||
				self.tokenizer.Symbol() == "/" ||
				self.tokenizer.Symbol() == "&amp;" ||
				self.tokenizer.Symbol() == "|" ||
				self.tokenizer.Symbol() == "&lt;" ||
				self.tokenizer.Symbol() == "&gt;" ||
				self.tokenizer.Symbol() == "=") {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // 'op'
			self.tokenizer.Advance()

			self.compileTerm()
		} else {
			break
		}
	}

	self.outputToXmlFile("</expression>")
}

func (self CompilationEngine) compileTerm() {
	self.outputToXmlFile("<term>")

	if self.tokenizer.TokenType() == "INT_CONST" {
		self.outputToXmlFile("<integerConstant> " + strconv.Itoa(self.tokenizer.IntVal()) + " </integerConstant>")
		self.tokenizer.Advance()
	} else if self.tokenizer.TokenType() == "STRING_CONST" {
		self.outputToXmlFile("<stringConstant> " + self.tokenizer.StringVal() + " </stringConstant>")
		self.tokenizer.Advance()
	} else if self.tokenizer.TokenType() == "KEYWORD" {
		self.outputToXmlFile("<keyword> " + self.tokenizer.KeyWord() + " </keyword>")
		self.tokenizer.Advance()
	} else if self.tokenizer.TokenType() == "IDENTIFIER" {
		name := self.tokenizer.Identifier()
		self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " used> " + name + " </identifier>") // varName
		self.tokenizer.Advance()

		if self.tokenizer.TokenType() == "SYMBOL" {
			if self.tokenizer.Symbol() == "[" { // '[' expression ']'
				self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '['
				self.tokenizer.Advance()

				self.compileExpression()

				self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ']'
				self.tokenizer.Advance()
			} else if self.tokenizer.Symbol() == "(" || self.tokenizer.Symbol() == "." {
				self.resolveSubroutineCall(true)
				self.tokenizer.Advance()
			}
		}
	} else if self.tokenizer.TokenType() == "SYMBOL" {
		if self.tokenizer.Symbol() == "-" || self.tokenizer.Symbol() == "~" {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '-|~'
			self.tokenizer.Advance()

			self.compileTerm()
		} else if self.tokenizer.Symbol() == "(" {
			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
			self.tokenizer.Advance()

			self.compileExpression()

			self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
			self.tokenizer.Advance()
		}
	}

	self.outputToXmlFile("</term>")
}

func (self CompilationEngine) compileExpressionList() {
	self.outputToXmlFile("<expressionList>")

	if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == ")" {
		// 何もしない
	} else {
		self.compileExpression()
		for {
			if self.tokenizer.TokenType() == "SYMBOL" && self.tokenizer.Symbol() == "," {
				self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ','
				self.tokenizer.Advance()

				self.compileExpression()
			} else {
				break
			}
		}
	}

	self.outputToXmlFile("</expressionList>")
}

func (self CompilationEngine) resolveType() string {
	typeName := ""

	if self.tokenizer.TokenType() == "KEYWORD" {
		typeName = self.tokenizer.KeyWord()
		self.outputToXmlFile("<keyword> " + typeName + " </keyword>") // ('int'|'char'|'boolean')
	} else if self.tokenizer.TokenType() == "IDENTIFIER" {
		name := self.tokenizer.Identifier()
		self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " used> " + name + " </identifier>") // className
	}

	return typeName
}

func (self CompilationEngine) resolveSubroutineCall(isCalledInTerm bool) {
	if isCalledInTerm {
		// 何もしない
	} else {
		name := self.tokenizer.Identifier()
		self.outputToXmlFile("<identifier category=" + self.symbolTable.KindOf(name) + " kind=" + self.symbolTable.KindOf(name) + " index=" + strconv.Itoa(self.symbolTable.IndexOf(name)) + " used> " + name + " </identifier>") // className
		self.tokenizer.Advance()
	}

	if self.tokenizer.Symbol() == "(" {
		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
		self.tokenizer.Advance()

		self.compileExpressionList()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
	} else if self.tokenizer.Symbol() == "." {
		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
		self.tokenizer.Advance()

		self.outputToXmlFile("<identifier> " + self.tokenizer.Identifier() + " </identifier>") // subroutineName
		self.tokenizer.Advance()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // '('
		self.tokenizer.Advance()

		self.compileExpressionList()

		self.outputToXmlFile("<symbol> " + self.tokenizer.Symbol() + " </symbol>") // ')'
	}
}

func (self CompilationEngine) isType() bool {
	if self.tokenizer.TokenType() == "KEYWORD" &&
		(self.tokenizer.KeyWord() == "int" || self.tokenizer.KeyWord() == "char" || self.tokenizer.KeyWord() == "boolean") {
		return true
	} else if self.tokenizer.TokenType() == "IDENTIFIER" {
		return true
	}

	return false
}

func (self CompilationEngine) isStatement() bool {
	if self.tokenizer.TokenType() != "KEYWORD" {
		return false
	}

	if self.tokenizer.KeyWord() == "let" || self.tokenizer.KeyWord() == "if" || self.tokenizer.KeyWord() == "while" || self.tokenizer.KeyWord() == "do" || self.tokenizer.KeyWord() == "return" {
		return true
	}

	return false
}

func (self CompilationEngine) outputToXmlFile(s string) {
	fmt.Fprintln(self.outputXmlFile, s)
}
