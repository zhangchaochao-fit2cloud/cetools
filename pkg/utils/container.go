package utils

import (
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"io"
)

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

	//preSystemCpuUsage := preCpuStats["system_cpu_usage"].(float64)
	preSystemCpuUsage := float64(0)
	if preCpuStats["system_cpu_usage"] != nil {
		preSystemCpuUsage = preCpuStats["system_cpu_usage"].(float64)
	}

	//systemCpu := cpuStatInfo["system_cpu_usage"].(float64)
	systemCpu := float64(0)
	if cpuStatInfo["system_cpu_usage"] != nil {
		systemCpu = cpuStatInfo["system_cpu_usage"].(float64)
	}

	cpuUsage := cpuStatInfo["cpu_usage"].(map[string]interface{})
	totalUsage := cpuUsage["total_usage"].(float64)

	//memLimit := SpaceFloat(memStatInfo["limit"].(float64))
	var memLimit string
	if memStatInfo["limit"] != nil {
		memLimit = SpaceFloat(memStatInfo["limit"].(float64))
	}

	//memUsed := SpaceFloat(memStatInfo["usage"].(float64))
	var memUsed string
	if memStatInfo["usage"] != nil {
		memUsed = SpaceFloat(memStatInfo["usage"].(float64))
	}

	cpuDelta := totalUsage - preTotalUsage
	systemCpuDelta := systemCpu - preSystemCpuUsage
	numberOfCores := cpuStatInfo["online_cpus"].(float64)
	var cpuPercent = ""
	if systemCpuDelta >= 0 {
		cpuPercent = Percent((cpuDelta / systemCpuDelta) * numberOfCores * 100)
	}
	return memUsed, memLimit, cpuPercent
}
