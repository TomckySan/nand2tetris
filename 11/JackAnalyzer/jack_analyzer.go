package main

import (
	"./modules"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getFileName(filePath string) string {
	splited := strings.Split(filePath, "/")
	return splited[len(splited)-1]
}

func main() {
	// symbolTable := modules.NewSymbolTable()

	// symbolTable.Define("foo", "int", "VAR")
	// symbolTable.Define("bar", "int", "STATIC")
	// symbolTable.Define("hoge", "int", "STATIC")
	// symbolTable.Define("piyo", "int", "ARG")
	// symbolTable.Define("fuga", "int", "VAR")
	// symbolTable.Define("x", "int", "STATIC")

	// fmt.Println(symbolTable.VarCount("STATIC"))
	// fmt.Println(symbolTable.VarCount("ARG"))
	// fmt.Println(symbolTable.VarCount("VAR"))

	// fmt.Println(symbolTable.KindOf("foo"))
	// fmt.Println(symbolTable.TypeOf("x"))
	// fmt.Println(symbolTable.IndexOf("bar"))
	// fmt.Println(symbolTable.IndexOf("x"))
	// os.Exit(0)

	source := os.Args[1]
	_ = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(filePath, ".jack") {
			return nil
		}
		tokenizer := modules.NewJackTokenizer(filePath)
		tokenizer.Advance()

		outputXmlFileName := strings.TrimSuffix(getFileName(filePath), ".jack") + ".xml"
		outputXmlFile, outputXmlFileErr := os.Create(outputXmlFileName)
		if outputXmlFileErr != nil {
			fmt.Println(outputXmlFileErr)
		}

		symbolTable := modules.NewSymbolTable()
		compilationEngine := modules.NewCompilationEngine(tokenizer, symbolTable, outputXmlFile)
		compilationEngine.CompileClass()

		return nil
	})
}
