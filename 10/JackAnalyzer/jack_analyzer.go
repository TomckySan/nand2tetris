package main

import (
	"./modules"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

		outputTXmlFileName := strings.TrimSuffix(getFileName(filePath), ".jack") + "T.xml"
		outputTXmlFile, outputTXmlFileErr := os.Create(outputTXmlFileName)
		if outputTXmlFileErr != nil {
			fmt.Println(outputTXmlFileErr)
		}

		fmt.Fprintln(outputTXmlFile, "<tokens>")
		for {
			if !tokenizer.Advance() {
				break
			}
			s := ""
			if tokenizer.TokenType() == "KEYWORD" {
				s = "<keyword> " + tokenizer.KeyWord() + " </keyword>"
			} else if tokenizer.TokenType() == "SYMBOL" {
				s = "<symbol> " + tokenizer.Symbol() + " </symbol>"
			} else if tokenizer.TokenType() == "IDENTIFIER" {
				s = "<identifier> " + tokenizer.Identifier() + " </identifier>"
			} else if tokenizer.TokenType() == "INT_CONST" {
				s = "<integerConstant> " + strconv.Itoa(tokenizer.IntVal()) + " </integerConstant>"
			} else if tokenizer.TokenType() == "STRING_CONST" {
				s = "<stringConstant> " + tokenizer.StringVal() + " </stringConstant>"
			}
			fmt.Fprintln(outputTXmlFile, s)
		}
		fmt.Fprintln(outputTXmlFile, "</tokens>")
		return nil
	})
}
