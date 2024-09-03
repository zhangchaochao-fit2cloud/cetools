package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"inspect/pkg/global"
	"inspect/pkg/service"
	"inspect/pkg/utils"
)

var inspect = &cobra.Command{
	Use:   "inspect",
	Short: "CloudExplorer 巡检",
	//Short: "CloudExplorer inspection",
	//	Long: `Generates shell completion code.
	//
	//Auto completion supports bash, zsh and fish. Output is to STDOUT.
	//
	//source <(kompose completion bash)
	//source <(kompose completion zsh)
	//kompose completion fish | source
	//
	//Will load the shell completion code.
	//	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		ss := service.NewISystemService()
		info, err := ss.LoadOsInfo()
		if err != nil {
			return err
		}
		ssInfo := ss.LoadCurrentInfo("all", "all")

		//global.Print.Title("**********%s**********", "系统信息")
		global.Print.Info("******** %s ********", "系统信息")
		fmt.Println("操作系统：", info.OS)
		fmt.Println("系统架构：", info.KernelArch)
		fmt.Println("内核版本：", info.KernelVersion)

		global.Print.Info("******** %s ********", "系统状态")
		//global.Print.Info("********    %s   ********", "系统状态")
		//global.Print.Title("**********%s**********", "CPU")
		global.Print.Info("内存")
		fmt.Println("内存：", utils.SpaceDisplay(ssInfo.MemoryTotal))
		fmt.Println("内存使用：", utils.SpaceDisplay(ssInfo.MemoryUsed))
		fmt.Println("内存剩余：", utils.SpaceDisplay(ssInfo.MemoryAvailable))
		fmt.Println("内存使用率：", ssInfo.CPUUsedPercent)

		global.Print.Info("交换内存")
		fmt.Println("交换内存：", utils.SpaceDisplay(ssInfo.SwapMemoryTotal))
		fmt.Println("交换内存使用：", utils.SpaceDisplay(ssInfo.SwapMemoryAvailable))
		fmt.Println("交换内存剩余：", utils.SpaceDisplay(ssInfo.SwapMemoryUsed))
		fmt.Println("交换内存使用率：", ssInfo.SwapMemoryUsedPercent)

		global.Print.Info("CPU")
		fmt.Println("CPU：", ssInfo.CPUTotal)
		fmt.Println("CPU使用：", ssInfo.CPUUsed)
		fmt.Println("CPUPercent：", ssInfo.CPUPercent)
		fmt.Println("CPU使用率：", ssInfo.CPUUsedPercent)
		fmt.Println("LoadUsagePercent：", ssInfo.LoadUsagePercent)
		fmt.Println("LoadUsagePercent：", ssInfo.LoadUsagePercent)

		global.Print.Info("IO")
		fmt.Println("IOReadBytes：", ssInfo.IOReadBytes)
		fmt.Println("IOWriteBytes：", ssInfo.IOWriteBytes)
		fmt.Println("IOCount：", ssInfo.IOCount)
		fmt.Println("IOReadTime：", ssInfo.IOReadTime)
		fmt.Println("IOWriteTime：", ssInfo.IOWriteTime)

		global.Print.Info("网络")
		fmt.Println("NetBytesSent：", ssInfo.NetBytesSent)
		fmt.Println("NetBytesRecv：", ssInfo.NetBytesRecv)

		return nil
	},
}

// Generate the appropriate autocompletion file
//func Generate(cmd *cobra.Command, args []string) error {
//	// Check the passed in arguments
//	if len(args) == 0 {
//		return fmt.Errorf("shell not specified. ex. kompose completion [bash|zsh|fish]")
//	}
//	if len(args) > 1 {
//		return fmt.Errorf("too many arguments. Expected only the shell type. ex. kompose completion [bash|zsh|fish]")
//	}
//
//	// Generate bash through cobra if selected
//	switch args[0] {
//	case "bash":
//		return cmd.Root().GenBashCompletion(os.Stdout)
//	case "zsh":
//		return runCompletionZsh(os.Stdout, cmd.Root())
//	case "fish":
//		return runCompletionFish(os.Stdout, cmd.Root())
//	default:
//		return fmt.Errorf("not a compatible shell, bash, zsh and fish are only supported")
//	}
//}

func init() {
	RootCmd.AddCommand(inspect)
}
