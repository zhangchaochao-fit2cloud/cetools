package service

import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"inspect/pkg/constant"
	"inspect/pkg/dto"
	"inspect/pkg/global"
	"inspect/pkg/utils"
	"inspect/pkg/utils/cmd"
	network "net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type diskInfo struct {
	Type   string
	Mount  string
	Device string
}

type SystemService struct{}

func NewISystemService() ISystemService {
	return &SystemService{}
}

type ISystemService interface {
	LoadMachineInfo() (*dto.MachineInfo, error)
	LoadBaseInfo(ioOption string, netOption string) (*dto.SystemBaseInfo, error)
	LoadCurrentInfo(ioOption string, netOption string) *dto.SystemCurrentInfo
}

func (s *SystemService) LoadMachineInfo() (*dto.MachineInfo, error) {
	var machineInfo dto.MachineInfo
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	machineInfo.Hostname = hostInfo.Hostname
	machineInfo.OS = hostInfo.OS
	machineInfo.Platform = hostInfo.Platform
	machineInfo.PlatformFamily = hostInfo.PlatformFamily
	machineInfo.PlatformVersion = hostInfo.PlatformVersion
	machineInfo.KernelArch = hostInfo.KernelArch
	machineInfo.KernelVersion = hostInfo.KernelVersion
	machineInfo.IpV4Addr = GetOutboundIP()
	machineInfo.OSInfo = GetOSInfoIP()

	cpuInfo, err := cpu.Info()
	if err == nil {
		machineInfo.CPUModelName = cpuInfo[0].ModelName
	}

	machineInfo.CPUCores, _ = cpu.Counts(false)
	machineInfo.CPULogicalCores, _ = cpu.Counts(true)

	//diskInfo, err := disk.Usage(global.CONF.System.BaseDir)
	//if err == nil {
	//	machineInfo.DiskSize = int64(diskInfo.Free)
	//}

	if machineInfo.KernelArch == "armv7l" {
		machineInfo.KernelArch = "armv7"
	}
	if machineInfo.KernelArch == "x86_64" {
		machineInfo.KernelArch = "amd64"
	}
	return &machineInfo, nil
}

func (u *SystemService) LoadBaseInfo(ioOption string, netOption string) (*dto.SystemBaseInfo, error) {
	var baseInfo dto.SystemBaseInfo
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	baseInfo.Hostname = hostInfo.Hostname
	baseInfo.OS = hostInfo.OS
	baseInfo.Platform = hostInfo.Platform
	baseInfo.PlatformFamily = hostInfo.PlatformFamily
	baseInfo.PlatformVersion = hostInfo.PlatformVersion
	baseInfo.KernelArch = hostInfo.KernelArch
	baseInfo.KernelVersion = hostInfo.KernelVersion
	ss, _ := json.Marshal(hostInfo)
	baseInfo.VirtualizationSystem = string(ss)

	//appInstall, err := appInstallRepo.ListBy()
	//if err != nil {
	//	return nil, err
	//}
	//baseInfo.AppInstalledNumber = len(appInstall)
	//postgresqlDbs, err := postgresqlRepo.List()
	//if err != nil {
	//	return nil, err
	//}
	//mysqlDbs, err := mysqlRepo.List()
	//if err != nil {
	//	return nil, err
	//}
	//baseInfo.DatabaseNumber = len(mysqlDbs) + len(postgresqlDbs)
	//website, err := websiteRepo.GetBy()
	//if err != nil {
	//	return nil, err
	//}
	//baseInfo.WebsiteNumber = len(website)
	//cronjobs, err := cronjobRepo.List()
	//if err != nil {
	//	return nil, err
	//}
	//baseInfo.CronjobNumber = len(cronjobs)

	cpuInfo, err := cpu.Info()
	if err == nil {
		baseInfo.CPUModelName = cpuInfo[0].ModelName
	}

	baseInfo.CPUCores, _ = cpu.Counts(false)
	baseInfo.CPULogicalCores, _ = cpu.Counts(true)

	//baseInfo.CurrentInfo = *u.LoadCurrentInfo(ioOption, netOption)
	return &baseInfo, nil
}

func (u *SystemService) LoadCurrentInfo(ioOption string, netOption string) *dto.SystemCurrentInfo {
	var currentInfo dto.SystemCurrentInfo
	hostInfo, _ := host.Info()
	currentInfo.Uptime = hostInfo.Uptime
	currentInfo.TimeSinceUptime = time.Now().Add(-time.Duration(hostInfo.Uptime) * time.Second).Format(constant.DateTimeLayout)
	currentInfo.Procs = hostInfo.Procs

	currentInfo.CPUTotal, _ = cpu.Counts(true)
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) == 1 {
		currentInfo.CPUUsedPercent = totalPercent[0]
		currentInfo.CPUUsed = currentInfo.CPUUsedPercent * 0.01 * float64(currentInfo.CPUTotal)
	}

	// 返回实时的所有 CPU 的使用率
	currentInfo.CPUPercent, _ = cpu.Percent(0, true)

	loadInfo, _ := load.Avg()
	currentInfo.Load1 = loadInfo.Load1
	currentInfo.Load5 = loadInfo.Load5
	currentInfo.Load15 = loadInfo.Load15
	currentInfo.LoadUsagePercent = loadInfo.Load1 / (float64(currentInfo.CPUTotal*2) * 0.75) * 100

	memoryInfo, _ := mem.VirtualMemory()
	currentInfo.MemoryTotal = memoryInfo.Total
	currentInfo.MemoryAvailable = memoryInfo.Available
	currentInfo.MemoryUsed = memoryInfo.Used
	currentInfo.MemoryUsedPercent = memoryInfo.UsedPercent

	swapInfo, _ := mem.SwapMemory()
	currentInfo.SwapMemoryTotal = swapInfo.Total
	currentInfo.SwapMemoryAvailable = swapInfo.Free
	currentInfo.SwapMemoryUsed = swapInfo.Used
	currentInfo.SwapMemoryUsedPercent = swapInfo.UsedPercent

	if hostInfo.OS != "darwin" {
		currentInfo.DiskData = loadDiskInfo()
		if ioOption == "all" {
			diskInfo, _ := disk.IOCounters()
			for _, state := range diskInfo {
				currentInfo.IOReadBytes += state.ReadBytes
				currentInfo.IOWriteBytes += state.WriteBytes
				currentInfo.IOCount += (state.ReadCount + state.WriteCount)
				currentInfo.IOReadTime += state.ReadTime
				currentInfo.IOWriteTime += state.WriteTime
			}
		} else {
			diskInfo, _ := disk.IOCounters(ioOption)
			for _, state := range diskInfo {
				currentInfo.IOReadBytes += state.ReadBytes
				currentInfo.IOWriteBytes += state.WriteBytes
				currentInfo.IOCount += (state.ReadCount + state.WriteCount)
				currentInfo.IOReadTime += state.ReadTime
				currentInfo.IOWriteTime += state.WriteTime
			}
		}

		if netOption == "all" {
			netInfo, _ := net.IOCounters(false)
			if len(netInfo) != 0 {
				currentInfo.NetBytesSent = netInfo[0].BytesSent
				currentInfo.NetBytesRecv = netInfo[0].BytesRecv
			}
		} else {
			netInfo, _ := net.IOCounters(true)
			for _, state := range netInfo {
				if state.Name == netOption {
					currentInfo.NetBytesSent = state.BytesSent
					currentInfo.NetBytesRecv = state.BytesRecv
				}
			}
		}

		currentInfo.ShotTime = time.Now()
	}
	//currentInfo.GPUData = loadGPUInfo()

	return &currentInfo
}

func loadDiskInfo() []dto.DiskInfo {
	var datas []dto.DiskInfo
	stdout, err := cmd.ExecWithTimeOut("df -hT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 2*time.Second)
	if err != nil {
		stdout, err = cmd.ExecWithTimeOut("df -lhT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 1*time.Second)
		if err != nil {
			return datas
		}
	}
	lines := strings.Split(stdout, "\n")

	var mounts []diskInfo
	var excludes = []string{"/mnt/cdrom", "/boot", "/boot/efi", "/dev", "/dev/shm", "/run/lock", "/run", "/run/shm", "/run/user"}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		if strings.HasPrefix(fields[6], "/snap") || len(strings.Split(fields[6], "/")) > 10 {
			continue
		}
		if strings.TrimSpace(fields[1]) == "tmpfs" {
			continue
		}
		if strings.Contains(fields[2], "K") {
			continue
		}
		if strings.Contains(fields[6], "docker") {
			continue
		}
		isExclude := false
		for _, exclude := range excludes {
			if exclude == fields[6] {
				isExclude = true
			}
		}
		if isExclude {
			continue
		}
		mounts = append(mounts, diskInfo{Type: fields[1], Device: fields[0], Mount: strings.Join(fields[6:], " ")})
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(mounts))
	for i := 0; i < len(mounts); i++ {
		go func(timeoutCh <-chan time.Time, mount diskInfo) {
			defer wg.Done()

			var itemData dto.DiskInfo
			itemData.Path = mount.Mount
			itemData.Type = mount.Type
			itemData.Device = mount.Device
			select {
			case <-timeoutCh:
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
				global.LOG.Errorf("load disk info from %s failed, err: timeout", mount.Mount)
			default:
				state, err := disk.Usage(mount.Mount)
				if err != nil {
					mu.Lock()
					datas = append(datas, itemData)
					mu.Unlock()
					global.LOG.Errorf("load disk info from %s failed, err: %v", mount.Mount, err)
					return
				}
				itemData.Total = state.Total
				itemData.Free = state.Free
				itemData.Used = state.Used
				itemData.UsedPercent = state.UsedPercent
				itemData.InodesTotal = state.InodesTotal
				itemData.InodesUsed = state.InodesUsed
				itemData.InodesFree = state.InodesFree
				itemData.InodesUsedPercent = state.InodesUsedPercent
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
			}
		}(time.After(5*time.Second), mounts[i])
	}
	wg.Wait()

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Path < datas[j].Path
	})
	return datas
}

func GetOutboundIP() string {
	conn, err := network.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return "IPNotFound"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*network.UDPAddr)
	return localAddr.IP.String()
}

func GetOSInfoIP() string {
	files := []string{"centos-release", "kylin-release", "redhat-release"}
	for _, file := range files {
		_, err := os.Stat("/etc/" + file)
		if os.IsNotExist(err) {
			continue
		}
		out, _ := cmd.ExecCmd("cat /etc/" + file)
		return utils.FilterSpecialChar(out)
	}
	if cmd.Exists("lsb_release") {
		out, _ := cmd.ExecCmd("lsb_release -r")
		return utils.FilterSpecialChar(out)
	}
	return ""
}

//func loadGPUInfo() []dto.GPUInfo {
//	list := xpack.LoadGpuInfo()
//	if len(list) == 0 {
//		return nil
//	}
//	var data []dto.GPUInfo
//	for _, gpu := range list {
//		var dataItem dto.GPUInfo
//		if err := common.Copy(&dataItem, &gpu); err != nil {
//			continue
//		}
//		dataItem.PowerUsage = dataItem.PowerDraw + " / " + dataItem.MaxPowerLimit
//		dataItem.MemoryUsage = dataItem.MemUsed + " / " + dataItem.MemTotal
//		data = append(data, dataItem)
//	}
//	return data
//}
