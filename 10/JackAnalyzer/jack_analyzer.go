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
	source := os.Args[1]
	_ = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(filePath, ".jack") {
			return nil
		}
		tokenizer := modules.NewJackTokenizer(filePath)
		tokenizer.Advance()

		// TXML
		// outputTXmlFileName := strings.TrimSuffix(getFileName(filePath), ".jack") + "T.xml"
		// outputTXmlFile, outputTXmlFileErr := os.Create(outputTXmlFileName)
		// if outputTXmlFileErr != nil {
		// 	fmt.Println(outputTXmlFileErr)
		// }
		// XML
		outputXmlFileName := strings.TrimSuffix(getFileName(filePath), ".jack") + ".xml"
		outputXmlFile, outputXmlFileErr := os.Create(outputXmlFileName)
		if outputXmlFileErr != nil {
			fmt.Println(outputXmlFileErr)
		}

		compilationEngine := modules.NewCompilationEngine(tokenizer, outputXmlFile)
		compilationEngine.CompileClass()

		return nil
	})
}
