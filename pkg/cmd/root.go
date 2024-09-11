package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"inspect/pkg/global"
	"inspect/pkg/init/log"
	"inspect/pkg/init/viper"
	"runtime"
	"strings"
)

var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var RootCmd = &cobra.Command{
	Use:           "ce-tool",
	Short:         "CloudExplorer 的巡检和高可用部署工具",
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Init()
		viper.Init()
		saveCommand(cmd)
	},
	Version: version(),
}

// Execute executes the root level command.
// It returns an error if any.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&global.Conf.Verbose, "verbose", "v", false, "详细输出")
	RootCmd.PersistentFlags().BoolVar(&global.Conf.SuppressWarnings, "suppress-warnings", false, "禁止所有异常")
	RootCmd.PersistentFlags().BoolVar(&global.Conf.ErrorOnWarning, "error-on-warning", false, "将任何警告视为错误")
	RootCmd.PersistentFlags().StringSliceVarP(&global.Conf.Files, "file", "f", []string{}, "配置文件")
	RootCmd.SetVersionTemplate(`{{.Version}}`)
}

func saveCommand(cmd *cobra.Command) {
	var cmdStr strings.Builder
	//cmdStr.WriteString("ce-tool inspect -r ")
	//var (
	//	wg sync.WaitGroup
	//)
	//wg.Add(len(cmd.Flags().Args()))
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flagPrefix := "-"
		if len(f.Name) > 1 {
			flagPrefix = "--"
		}
		if f.Name == "remote" || f.Name == "r" || f.Name == "file" {
			return
		}
		value := f.Value.String()
		if f.Value.Type() == "bool" {
			if f.Value.String() == "false" {
				return
			}
			value = ""
			return
		}
		cmdStr.WriteString(fmt.Sprintf(" %s%s %v", flagPrefix, f.Name, value))
		//wg.Done()
	})
	//wg.Wait()
	global.Conf.Command = cmdStr.String()
}

func version() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Version:    %s\n", Version))
	sb.WriteString(fmt.Sprintf("Git commit: %s\n", GitCommit))
	sb.WriteString(fmt.Sprintf("Build time: %s\n", BuildTime))
	sb.WriteString(fmt.Sprintf("Go version: %s\n", runtime.Version()))
	sb.WriteString(fmt.Sprintf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH))
	return sb.String()
}
