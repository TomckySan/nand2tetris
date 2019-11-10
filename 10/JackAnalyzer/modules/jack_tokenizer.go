package modules

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type JackTokenizer struct {
	inputJackFile          *os.File
	inputJackFileReader    *bufio.Reader
	currentLineTokensIndex int
	currentLineTokens      []string
	currentToken           string
	isInComment            bool
}

func (self *JackTokenizer) initialize(inputJackFilePath string) {
	file, err := os.Open(inputJackFilePath)
	if err != nil {
		file.Close()
		fmt.Println(err)
		os.Exit(1)
	}
	self.inputJackFile = file
	self.inputJackFileReader = bufio.NewReaderSize(file, 1024)
	self.currentLineTokensIndex = 0
	self.currentLineTokens = []string{}
	self.currentToken = ""
	self.isInComment = false
}

func NewJackTokenizer(inputJackFilePath string) *JackTokenizer {
	jackTokenizer := &JackTokenizer{}
	jackTokenizer.initialize(inputJackFilePath)
	return jackTokenizer
}

func (self *JackTokenizer) Advance() bool {
	if self.currentLineTokensIndex > (len(self.currentLineTokens) - 1) {
		self.currentLineTokensIndex = 0
		self.currentLineTokens = []string{}
	}

	if self.currentLineTokensIndex == 0 {
		for len(self.currentLineTokens) == 0 {
			line, _, err := self.inputJackFileReader.ReadLine()
			if err == io.EOF {
				return false
			} else if err != nil {
				self.inputJackFile.Close()
				fmt.Println(err)
				os.Exit(1)
			}

			self.currentLineTokens = self.splitToTokens(self.removedCommentLine(string(line)))
		}
	}

	self.currentToken = self.currentLineTokens[self.currentLineTokensIndex]
	// fmt.Println(self.currentToken) // TODO: Advanceを呼ぶごとにcurrentTokenを確認（デバッグ）
	self.currentLineTokensIndex++

	return true
}

func (self *JackTokenizer) removedCommentLine(line string) string {
	containCommentEndToken := strings.Contains(line, "*/")
	if self.isInComment && !containCommentEndToken {
		// 複数行コメントの途中の行、かつ複数行コメントの終了が含まれてない場合は空文字を返す
		return ""
	}

	// 1行コメントが見つかったら末尾まで削除
	line = strings.Split(line, "//")[0]

	containCommentBeginToken := strings.Contains(line, "/*")
	if containCommentBeginToken && !containCommentEndToken {
		// 複数行コメントの開始が含まれており、かつ終了が含まれてなければ開始から末尾まで削除
		line = strings.Split(line, "/*")[0]
		self.isInComment = true
	} else if !containCommentBeginToken && containCommentEndToken {
		// 複数行コメントの開始が含まれておらず、かつ終了が含まれているなら終了まで削除し以降は残す
		line = strings.SplitN(line, "*/", 2)[1]
		self.isInComment = false
		self.removedCommentLine(line) // 終了が含まれなくなるまで再帰
	} else if containCommentBeginToken && containCommentEndToken {
		if strings.Index(line, "*/") < strings.Index(line, "/*") {
			line = strings.SplitAfterN(line, "*/", 2)[1]
		} else {
			s1 := strings.SplitN(line, "/*", 2)
			s2 := strings.SplitN(s1[1], "*/", 2)
			line = s1[0] + " " + s2[1]
		}
		self.isInComment = false
		self.removedCommentLine(line) // 終了が含まれなくなるまで再帰
	}

	return line
}

func (self JackTokenizer) splitToTokens(line string) []string {
	result := []string{}
	strExclusionSymbol := ""
	isInStringConstant := false
	isInSpace := false

	for _, r := range line {
		if strings.Contains(string(r), "\"") {
			// "であれば文字列の開始か終了のいずれか
			isInStringConstant = !isInStringConstant
			strExclusionSymbol += string(r)
		} else if !isInStringConstant && strings.ContainsAny(string(r), "{}()[].,;+-*/&|<>=~") {
			// 文字列の中に含まれていないシンボルであればappendする
			if strExclusionSymbol != "" {
				result = append(result, strExclusionSymbol)
			}
			strExclusionSymbol = ""
			result = append(result, string(r))
		} else if !isInStringConstant && unicode.IsSpace(r) {
			// 文字列の中に含まれていないスペースは除外するため準備
			isInSpace = true
		} else if !isInStringConstant && !unicode.IsSpace(r) && isInSpace {
			// スペースでなくなった場合appendして次の開始文字を設定
			isInSpace = false
			if strExclusionSymbol != "" {
				result = append(result, strExclusionSymbol)
			}
			strExclusionSymbol = string(r)
		} else {
			strExclusionSymbol += string(r)
		}
	}

	if strExclusionSymbol != "" {
		result = append(result, strExclusionSymbol)
	}

	return result
}

func (self JackTokenizer) TokenType() string {
	if self.isKeyWord() {
		return "KEYWORD"
	}
	if self.isSymbol() {
		return "SYMBOL"
	}
	if self.isIntConst() {
		return "INT_CONST"
	}
	if self.isStringConst() {
		return "STRING_CONST"
	}
	if self.isIdentifier() {
		return "IDENTIFIER"
	}
	return ""
}

func (self JackTokenizer) KeyWord() string {
	return self.currentToken
}

func (self JackTokenizer) Symbol() string {
	if self.currentToken == "<" {
		return "&lt;"
	} else if self.currentToken == ">" {
		return "&gt;"
	} else if self.currentToken == "&" {
		return "&amp;"
	} else {
		return self.currentToken
	}
}

func (self JackTokenizer) Identifier() string {
	return self.currentToken
}

func (self JackTokenizer) IntVal() int {
	i, _ := strconv.Atoi(self.currentToken)
	return i
}

func (self JackTokenizer) StringVal() string {
	return strings.Trim(self.currentToken, "\"")
}

func (self JackTokenizer) isKeyWord() bool {
	keyWords := []string{
		"class", "constructor", "function", "method", "field",
		"static", "var", "int", "char", "boolean", "void",
		"true", "false", "null", "this", "let", "do", "if",
		"else", "while", "return",
	}

	for _, keyWord := range keyWords {
		if self.currentToken == keyWord {
			return true
		}
	}

	return false
}

func (self JackTokenizer) isSymbol() bool {
	symbols := []string{
		"{", "}", "(", ")", "[", "]", ".",
		",", ";", "+", "-", "*", "/", "&",
		"|", "<", ">", "=", "~",
	}

	for _, symbol := range symbols {
		if self.currentToken == symbol {
			return true
		}
	}

	return false
}

func (self JackTokenizer) isIdentifier() bool {
	return regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z]*$`).MatchString(self.currentToken)
}

func (self JackTokenizer) isIntConst() bool {
	if regexp.MustCompile(`^[0-9]*$`).MatchString(self.currentToken) {
		i, _ := strconv.Atoi(self.currentToken)
		if 0 <= i && i <= 32767 {
			return true
		}
	}

	return false
}

func (self JackTokenizer) isStringConst() bool {
	return regexp.MustCompile(`^".*"$`).MatchString(self.currentToken)
}

func (self JackTokenizer) GetCurrentToken() string {
	return self.currentToken
}
