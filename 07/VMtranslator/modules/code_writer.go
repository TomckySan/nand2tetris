package modules

import (
	"fmt"
	"os"
	"strings"
)

type CodeWriter struct {
	filePath string
	asmFile  *os.File
	labelNum int
}

func (self *CodeWriter) initialize() {
	// Open output file
	vmFilePathTrimExtension := strings.TrimSuffix(self.filePath, ".vm")
	asmFileName := fmt.Sprintf("%s.asm", vmFilePathTrimExtension)
	asmFile, asmFileOpenErr := os.Create(asmFileName)
	if asmFileOpenErr != nil {
		fmt.Println(asmFileOpenErr)
	}

	self.asmFile = asmFile
	self.labelNum = 0
}

func NewCodeWriter(filePath string) *CodeWriter {
	codeWriter := &CodeWriter{filePath: filePath}
	codeWriter.initialize()
	return codeWriter
}

// func (self CodeWriter) setFileName(fileName string) {
// 	// vmFile := fileName
// }

func (self *CodeWriter) writeBinaryFunction(command string) {
	self.writePop()
	fmt.Fprintln(self.asmFile, "D=M") // Mの値をDレジスタへ
	self.writePop()
	switch command {
	case "add":
		fmt.Fprintln(self.asmFile, "D=D+M") // Mの値とDレジスタの値を加算してDレジスタへ上書き
	case "sub":
		fmt.Fprintln(self.asmFile, "D=M-D") // Mの値とDレジスタの値を減算してDレジスタへ上書き
	case "and":
		fmt.Fprintln(self.asmFile, "D=D&M") // Mの値とDレジスタの値のANDをとってDレジスタへ上書き
	case "or":
		fmt.Fprintln(self.asmFile, "D=D|M") // Mの値とDレジスタの値のORをとってDレジスタへ上書き
	}
	self.writePush() // 結果をpush
}

func (self *CodeWriter) writeBinaryCompareFunction(command string) {
	self.writePop()                   // 変数をpop
	fmt.Fprintln(self.asmFile, "D=M") // Mの値をDレジスタへ
	self.writePop()                   // 変数をpop
	label1 := self.getNewLabel()
	label2 := self.getNewLabel()
	commandType := ""
	switch command {
	case "eq":
		commandType = "JEQ"
	case "gt":
		commandType = "JGT"
	case "lt":
		commandType = "JLT"
	default:
		commandType = ""
	}
	fmt.Fprintln(self.asmFile, "D=M-D") // Mの値とDレジスタの値を減算してDレジスタへ上書き
	fmt.Fprintln(self.asmFile, fmt.Sprintf("@%s", label1))
	fmt.Fprintln(self.asmFile, fmt.Sprintf("D;%s", commandType)) // 条件一致するならlabel1へ
	fmt.Fprintln(self.asmFile, "D=0")                            // ↑でジャンプしないならfalseなので0を設定
	fmt.Fprintln(self.asmFile, fmt.Sprintf("@%s", label2))
	fmt.Fprintln(self.asmFile, "0;JMP")                     // 強制的にlabel2へ
	fmt.Fprintln(self.asmFile, fmt.Sprintf("(%s)", label1)) // 条件一致したらここへジャンプ
	fmt.Fprintln(self.asmFile, "D=-1")                      // ここにくるならtrueなので-1を設定
	fmt.Fprintln(self.asmFile, fmt.Sprintf("(%s)", label2)) // 条件一致しない場合は↑を通過せずにここへジャンプ
	self.writePush()                                        // 結果をpush
}

func (self *CodeWriter) getNewLabel() string {
	self.labelNum++
	return fmt.Sprintf("LABEL%d", self.labelNum)
}

func (self *CodeWriter) writeUnaryFunction(command string) {
	self.writePop()
	switch command {
	case "neg":
		fmt.Fprintln(self.asmFile, "D=-M") // Mの値を符号反転してDレジスタへ
	case "not":
		fmt.Fprintln(self.asmFile, "D=!M") // Mの値のNOTをとってDレジスタへ
	}
	self.writePush() // 結果をpush
}

func (self *CodeWriter) writePush() {
	fmt.Fprintln(self.asmFile, "@SP")   // SPのアドレスを参照
	fmt.Fprintln(self.asmFile, "A=M")   // SP値をAレジスタへ（MはSP値と同じアドレスを参照するようになる）
	fmt.Fprintln(self.asmFile, "M=D")   // MへDレジスタの値を設定
	fmt.Fprintln(self.asmFile, "@SP")   // SPのアドレスを参照
	fmt.Fprintln(self.asmFile, "M=M+1") // SP値を加算→SPが指す場所を1つ進める
}

func (self *CodeWriter) writePop() {
	fmt.Fprintln(self.asmFile, "@SP")   // SPのアドレスを参照
	fmt.Fprintln(self.asmFile, "M=M-1") // SP値を減算→SPが指す場所を1つ戻す
	fmt.Fprintln(self.asmFile, "A=M")   // SP値をAレジスタへ（MはSP値と同じアドレスを参照するようになる）
}

func (self *CodeWriter) WriteArithmetic(command string) {
	if command == "add" || command == "sub" || command == "and" || command == "or" {
		self.writeBinaryFunction(command)
	} else if command == "eq" || command == "gt" || command == "lt" {
		self.writeBinaryCompareFunction(command)
	} else if command == "neg" || command == "not" {
		self.writeUnaryFunction(command)
	}
}

func (self CodeWriter) WritePushPop(command string, segment string, index int) {
	// push constant index
	fmt.Fprintln(self.asmFile, fmt.Sprintf("@%d", index))
	fmt.Fprintln(self.asmFile, "D=A")
	fmt.Fprintln(self.asmFile, "@SP")   // SPのアドレスを参照
	fmt.Fprintln(self.asmFile, "A=M")   // SP値をAレジスタへ（MはSP値と同じアドレスを参照するようになる）
	fmt.Fprintln(self.asmFile, "M=D")   // MへDレジスタの値を設定
	fmt.Fprintln(self.asmFile, "@SP")   // SPのアドレスを参照
	fmt.Fprintln(self.asmFile, "M=M+1") // SP値を加算→SPが指す場所を1つ進める
}

func (self *CodeWriter) Close() {
	self.asmFile.Close()
}
