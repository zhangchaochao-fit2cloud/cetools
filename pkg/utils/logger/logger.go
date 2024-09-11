package logger

import (
	"cetool/pkg/constant"
	"fmt"
	"strings"
)

type Logger struct {
}

func (log *Logger) Format(prefix, format string, a ...any) {
	fmt.Print(log.SFormat(prefix, format, a...))
}

func (log *Logger) SFormat(prefix, format string, a ...any) string {
	var content string = format
	if len(a) > 0 {
		content = fmt.Sprintf(format, a...)
		content = strings.ReplaceAll(content, "[", "")
		content = strings.ReplaceAll(content, "]", "")
		content = strings.ReplaceAll(content, "%!(EXTRA interface {}=)", "")
	}
	return fmt.Sprintf("%s%s%s\n", prefix, content, constant.Reset)
}

func (log *Logger) Debug(format string, a ...any) {
	log.Format(constant.Red, format, a...)
}
func (log *Logger) SDebug(format string, a ...any) string {
	return log.SDebug(constant.Red, format, a)
}

func (log *Logger) Info(format string, a ...any) {
	log.Format(constant.Green, format, a...)
}

func (log *Logger) SInfo(format string, a ...any) string {
	return log.SFormat(constant.Green, format, a...)
}

func (log *Logger) Title(format string, a string) {
	// 定义总宽度为30的字符串
	width := 20
	// 计算两边空格的长度
	padding := width - len(a)
	// 计算左侧空格数
	leftPadding := padding / 2

	// 计算右侧空格数
	rightPadding := padding - leftPadding
	content := fmt.Sprintf("%s%s%s", strings.Repeat(" ", leftPadding), a, strings.Repeat(" ", rightPadding))
	log.Format(constant.Green, format, content)
}

func (log *Logger) Warning(format string, a ...any) {
	log.Format(constant.Yellow, format, a...)
}

func (log *Logger) SWarning(format string, a ...any) string {
	return log.SFormat(constant.Yellow, format, a...)
}

func (log *Logger) Error(format string, a ...any) {
	log.Format(constant.Red, format, a...)
}

func (log *Logger) SError(format string, a ...any) string {
	return log.SFormat(constant.Red, format, a...)
}

//func NewLogger() *Logger {
//	return &Logger{}
//}
