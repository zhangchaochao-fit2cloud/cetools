package utils

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"inspect/pkg/constant"
	"inspect/pkg/global"
	"io"
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

func SpaceFloat(sizeKB float64) string {
	switch {
	case sizeKB >= float64(constant.PB):
		return fmt.Sprintf("%.1fP", sizeKB/float64(constant.PB))
	case sizeKB >= float64(constant.TB):
		return fmt.Sprintf("%.1fT", sizeKB/float64(constant.TB))
	case sizeKB >= float64(constant.GB):
		return fmt.Sprintf("%.1fG", sizeKB/float64(constant.GB))
	case sizeKB >= float64(constant.MB):
		return fmt.Sprintf("%.1fM", sizeKB/float64(constant.MB))
	default:
		return fmt.Sprintf("%fK", sizeKB)
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

func Stat(stats container.StatsResponseReader) map[string]interface{} {
	data, err := io.ReadAll(stats.Body)
	if err != nil {
		panic(err)
	}

	var statsData map[string]interface{}
	if err := json.Unmarshal(data, &statsData); err != nil {
		panic(err)
	}
	return statsData
}

func GetContainerStats(statMap map[string]interface{}) (string, string, string) {
	memStatInfo := statMap["memory_stats"].(map[string]interface{})
	cpuStatInfo := statMap["cpu_stats"].(map[string]interface{})
	preCpuStats := statMap["precpu_stats"].(map[string]interface{})
	preTotalUsage := preCpuStats["cpu_usage"].(map[string]interface{})["total_usage"].(float64)
	preSystemCpuUsage := preCpuStats["system_cpu_usage"].(float64)

	systemCpu := cpuStatInfo["system_cpu_usage"].(float64)
	cpuUsage := cpuStatInfo["cpu_usage"].(map[string]interface{})
	totalUsage := cpuUsage["total_usage"].(float64)

	memLimit := SpaceFloat(memStatInfo["limit"].(float64))
	memUsed := SpaceFloat(memStatInfo["usage"].(float64))

	cpuDelta := totalUsage - preTotalUsage
	systemCpuDelta := systemCpu - preSystemCpuUsage
	numberOfCores := cpuStatInfo["online_cpus"].(float64)
	cpuPercent := Percent((cpuDelta / systemCpuDelta) * numberOfCores * 100)

	return memUsed, memLimit, cpuPercent
}
