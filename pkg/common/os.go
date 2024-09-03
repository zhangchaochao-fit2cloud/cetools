package common

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"strconv"
)

type LinuxSystemInfo struct {
	CPU int
	Mem uint64
	OS  string
}

func (s LinuxSystemInfo) PrintInfo() string {
	return fmt.Sprintf("CPU: %sC, Mem: %dG, 操作系统：%v", strconv.Itoa(s.CPU), s.Mem, s.OS)
}

func osSystem() string {
	switch runtime.GOOS {
	case "darwin":
		//osVersion := runShellCommand("sw_vers -productVersion")
		return ""
	case "linux":
		break
	default:
		break
	}
	return "未知"
}

func GetSystemInfo() *LinuxSystemInfo {

	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	var cores int
	for _, stat := range cpuInfo {
		cores += int(stat.Cores)
	}

	//osName := runtime.GOOS
	//fmt.Println("操作系统:", osName)

	//hostName, err := os.Hostname()
	//if err == nil {
	//	fmt.Println("主机名:", hostName)
	//}

	//var osVersion string
	//switch osName {
	//case "windows":
	//	osVersion = os.Getenv("OS")
	//case "linux":
	//	//releaseFile, err := os.ReadFile("/etc/os-release")
	//	//if err == nil {
	//	//	osVersion = parseOSRelease(string(releaseFile))
	//	//}
	//case "darwin":
	//	//osVersion = runShellCommand("sw_vers -productVersion")
	//default:
	//	osVersion = "未知版本"
	//}

	//fmt.Println("操作系统版本:", osVersion)
	//
	//fmt.Printf("操作系统: %d\n", osName)
	////fmt.Printf("磁盘总量: %d\n", diskInfo.Total/1024/1024/1024)
	////fmt.Printf("磁盘剩余: %d\n", diskInfo.Free/1024/1024/1024)
	//fmt.Printf("CPU 核数: %d\n", cores)
	//fmt.Printf("内存大小: %d GB\n", memInfo.Total/1024/1024/1024)
	var memoryTotal = memInfo.Total / 1024 / 1024 / 1024
	os := osSystem()
	return &LinuxSystemInfo{
		CPU: cores,
		Mem: memoryTotal,
		OS:  os,
	}
}
