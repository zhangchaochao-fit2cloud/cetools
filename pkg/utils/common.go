package utils

import (
	"fmt"
	"inspect/pkg/constant"
	"inspect/pkg/global"
	"strings"
)

func Space(sizeKB uint64) string {
	switch {
	case sizeKB >= uint64(constant.PB):
		return fmt.Sprintf("%.1fP", float64(sizeKB)/float64(constant.PB))
	case sizeKB >= uint64(constant.TB):
		return fmt.Sprintf("%.1fT", float64(sizeKB)/float64(constant.TB))
	case sizeKB >= uint64(constant.GB):
		return fmt.Sprintf("%.1fG", float64(sizeKB)/float64(constant.GB))
	case sizeKB >= uint64(constant.MB):
		return fmt.Sprintf("%.1fM", float64(sizeKB)/float64(constant.MB))
	default:
		return fmt.Sprintf("%dK", sizeKB)
	}
}

func Percent(percent float64) string {
	switch {
	case percent >= 80:
		return FilterStrSpecialChar(global.Print.SError("%.2f%%", percent))
	case percent >= 60:
		return FilterStrSpecialChar(global.Print.SWarning("%.2f%%", percent))
	default:
		//return FilterStrSpecialChar(global.Print.SInfo("%.2f%%", percent))
		return fmt.Sprintf("%.2f%%", percent)
	}
}

func FPercent(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
func FormatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func CalculatePercent(total float64, used float64) string {
	if total == 0 {
		return Percent(0)
	}
	return Percent((used / total) * 100)
}

func FilterSpecialChar(out []byte) string {
	return FilterStrSpecialChar(string(out))
}

func FilterStrSpecialChar(out string) string {
	fileContent := strings.ReplaceAll(out, "\r", "")
	fileContent = strings.ReplaceAll(fileContent, "\n", "")
	return fileContent
}
