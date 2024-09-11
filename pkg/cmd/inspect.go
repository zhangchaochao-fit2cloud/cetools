package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"inspect/pkg/global"
	"inspect/pkg/service"
	u "inspect/pkg/utils"
	"inspect/pkg/utils/cmp"
	"inspect/pkg/utils/table"
)

var (
	//GlobalService       bool
	GlobalShowInfo      string
	GlobalRemoteCommand bool
	GlobalUploadBin     bool
)

var Inspect = &cobra.Command{
	Use:   "inspect",
	Short: "CloudExplorer 巡检",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := GenerateLocalInfo()
		if err != nil {
			//panic(err)
			global.Print.Error(err.Error())
			return nil
		}
		if !GlobalRemoteCommand {
			err = GenerateNodeInfo()
			if err != nil {
				global.Print.Error(err.Error())
				return nil
			}
		}
		return nil
	},
}

func GenerateNodeInfo() error {
	if global.CONF == nil || len(global.CONF.Nodes) <= 0 {
		return nil
	}

	for _, node := range global.CONF.Nodes {
		if len(node.Addr) <= 0 {
			continue
		}
		global.Print.Info("**************** %s ****************", "节点 "+node.Addr+" 信息")
		ss := service.SshService

		command := "ce-tool inspect -n "
		if GlobalUploadBin {
			tools, err := ss.UploadTools(node)
			if err != nil {
				u.CheckError(tools, err)
				return nil
			}
			command = tools + " inspect -n "
		}

		run, err := ss.Run(node, command+global.Conf.Command)
		if err != nil {
			u.CheckError(run, err)
			return nil
		}
		fmt.Println(run)
	}

	return nil
}

// GenerateLocalInfo 打印本地的信息
func GenerateLocalInfo() error {
	if !GlobalRemoteCommand {
		global.Print.Info("**************** %s ****************", "本机信息")
	}
	ss := service.NewISystemService()

	var (
		machineInfo = false
		serviceInfo = false
	)
	if len(GlobalShowInfo) <= 0 || GlobalShowInfo == "all" {
		machineInfo = true
		serviceInfo = true
	}

	if GlobalShowInfo == "service" {
		serviceInfo = true
	}

	if GlobalShowInfo == "info" {
		machineInfo = true
	}

	if machineInfo {
		// 机器信息
		err := GenerateMachineInfo(ss)
		if err != nil {
			return err
		}

		// 操作系统信息
		GenerateSystemInfo(ss)
	}

	if serviceInfo {
		// 服务信息
		err := GenerateServiceInfo()
		return err
	}

	return nil
}
func GenerateMachineInfo(ss service.ISystemService) error {
	info, err := ss.LoadMachineInfo()
	if err != nil {
		return err
	}
	global.Print.Info("【 %s 】", "系统信息")
	var osInfo string
	if len(info.OSInfo) > 0 {
		osInfo = info.OSInfo
	} else {
		osInfo = info.OS
	}
	fmt.Println("  操作系统：", osInfo)
	fmt.Println("  内核版本：", info.KernelVersion)
	fmt.Println("  CPU 架构：", info.KernelArch)
	fmt.Println("  CPU 型号：", info.CPUModelName)
	fmt.Println("  主 机 名：", info.Hostname)
	fmt.Println("  IP  地址：", info.IpV4Addr)
	return nil
}

func GenerateServiceInfo() error {
	services, err := cmp.GetCmpServices()
	global.Print.Info("【 %s 】", "服务信息")
	if err != nil {
		return err
	}
	var data [][]string
	for _, ser := range services {
		item := []string{ser.Name, ser.Status, ser.Ports, ser.CPUPercent, ser.MemUsed, ser.MemLimit, ser.Runtime}
		data = append(data, item)
	}

	table.Print([]string{"服务名", "状态", "开放端口", "CPU 使用", "内存使用", "内存限制", "运行时间"}, data)
	return nil
}

func GenerateSystemInfo(ss service.ISystemService) {
	ssInfo := ss.LoadCurrentInfo("all", "all")
	global.Print.Info("【 %s 】", "系统状态")

	fmt.Printf("  负载：1分钟：%s, 5分钟：%s, 15分钟：%s\n", u.Percent(ssInfo.Load1), u.Percent(ssInfo.Load5), u.Percent(ssInfo.Load15))

	var diskTotal uint64
	var diskUsed uint64
	for _, disk := range ssInfo.DiskData {
		diskTotal += disk.Total
		diskUsed += disk.Used
	}
	data := [][]string{
		{"CPU", u.FormatFloat(float64(ssInfo.CPUTotal)), u.FormatFloat(float64(ssInfo.CPUTotal) - ssInfo.CPUUsed), u.Percent(ssInfo.CPUUsedPercent)},
		{"内存", u.Space(ssInfo.MemoryTotal), u.Space(ssInfo.MemoryTotal - ssInfo.MemoryUsed), u.Percent(ssInfo.MemoryUsedPercent)},
		{"交换内存", u.Space(ssInfo.SwapMemoryTotal), u.Space(ssInfo.SwapMemoryAvailable), u.Percent(ssInfo.SwapMemoryUsedPercent)},
		//{"磁盘", u.Space(diskTotal), u.Space(diskUsed), u.CalculatePercent(float64(diskTotal), float64(diskUsed))},
	}
	table.Print([]string{"资源", "总量", "剩余", "使用率"}, data)

	if len(ssInfo.DiskData) > 0 {
		global.Print.Info("【 %s 】", "磁盘状态")
		var spaceData [][]string
		for _, disk := range ssInfo.DiskData {
			item := []string{disk.Path, u.Space(disk.Total), u.Space(disk.Free), u.Percent(disk.UsedPercent)}
			spaceData = append(spaceData, item)
		}
		table.Print([]string{"挂载点", "总量", "剩余", "使用率"}, spaceData)
	}
}

func init() {
	Inspect.PersistentFlags().BoolVarP(&GlobalRemoteCommand, "is-node", "n", false, "is Remote Command ")
	if f := Inspect.PersistentFlags().Lookup("is-node"); f != nil {
		f.Hidden = true
	}

	Inspect.PersistentFlags().BoolVarP(&GlobalUploadBin, "upload-tool", "p", false, "自动上传工具到节点")
	//Inspect.PersistentFlags().BoolVarP(&GlobalService, "service", "s", false, "输出服务状态")
	Inspect.PersistentFlags().StringVarP(&GlobalShowInfo, "info", "i", "all", "查看巡检：all 全部，service 服务，info 主机信息")
	RootCmd.AddCommand(Inspect)
}
