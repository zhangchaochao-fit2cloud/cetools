package common

import (
	"fmt"
	"inspect/pkg/constant"
	"strings"
)

type Logger struct {
}

func (log *Logger) Format(prefix, format string, a ...any) {
	var content string = format
	if len(a) > 0 {
		content = fmt.Sprintf(format, a...)
	}
	fmt.Printf("%s%s%s\n", prefix, content, constant.Reset)
}

func (log *Logger) Debug(format string, a ...any) {
	log.Format(constant.Red, format, a...)
}

func (log *Logger) Info(format string, a ...any) {
	log.Format(constant.Green, format, a...)
}

func (log *Logger) Title(format string, a string) {
	// 定义总宽度为30的字符串
	width := 20
	// 计算两边空格的长度
	padding := width - len(a)
	//padding := width - utf8.RuneCountInString(a)
	//fmt.Println("length", utf8.RuneCountInString(a))
	// 计算左侧空格数
	leftPadding := padding / 2

	// 计算右侧空格数
	rightPadding := padding - leftPadding
	//fmt.Println("leftPadding", leftPadding)
	//fmt.Println("rightPadding", rightPadding)
	content := fmt.Sprintf("%s%s%s", strings.Repeat(" ", leftPadding), a, strings.Repeat(" ", rightPadding))
	log.Format(constant.Green, format, content)
}

func (log *Logger) Warning(format string, a ...any) {
	log.Format("WARNING", format, a...)
}

func (log *Logger) Error(format string, a ...any) {
	log.Format(constant.Red, format, a...)
}

func NewLogger() *Logger {
	return &Logger{}
}
