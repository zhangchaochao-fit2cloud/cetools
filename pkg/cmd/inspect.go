package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"inspect/pkg/global"
	"inspect/pkg/service"
	u "inspect/pkg/utils"
	"os"
)

var inspect = &cobra.Command{
	Use:   "inspect",
	Short: "CloudExplorer 巡检",
	RunE: func(cmd *cobra.Command, args []string) error {
		ss := service.NewISystemService()

		// 机器信息
		err := GenerateMachineInfo(ss)
		if err != nil {
			return err
		}

		// 操作系统信息
		GenerateSystemInfo(ss)

		return nil
	},
}

func GenerateMachineInfo(ss service.ISystemService) error {
	info, err := ss.LoadMachineInfo()
	if err != nil {
		return err
	}
	global.Print.Info("**************** %s ****************", "系统信息")
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

func GenerateServiceInfo(ss service.ISystemService) {

}
func GenerateSystemInfo(ss service.ISystemService) {
	ssInfo := ss.LoadCurrentInfo("all", "all")
	global.Print.Info("**************** %s ****************", "系统状态")

	fmt.Printf("  负载：1分钟：%s, 5分钟：%s, 15分钟：%s\n", u.Percent(ssInfo.Load1), u.Percent(ssInfo.Load5), u.Percent(ssInfo.Load15))

	var diskTotal uint64
	var diskUsed uint64
	for _, disk := range ssInfo.DiskData {
		diskTotal += disk.Total
		diskUsed += disk.Used
	}
	data := [][]string{
		{"CPU", u.FormatFloat(float64(ssInfo.CPUTotal)), u.FormatFloat(ssInfo.CPUUsed), u.Percent(ssInfo.CPUUsedPercent)},
		{"内存", u.Space(ssInfo.MemoryTotal), u.Space(ssInfo.MemoryUsed), u.Percent(ssInfo.MemoryUsedPercent)},
		{"交换内存", u.Space(ssInfo.SwapMemoryTotal), u.Space(ssInfo.SwapMemoryAvailable), u.Percent(ssInfo.SwapMemoryUsedPercent)},
		{"磁盘", u.Space(diskTotal), u.Space(diskUsed), u.CalculatePercent(float64(diskTotal), float64(diskUsed))},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"资源", "总量", "使用", "使用率"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func init() {
	RootCmd.AddCommand(inspect)
}
