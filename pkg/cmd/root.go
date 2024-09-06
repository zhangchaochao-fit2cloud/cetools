package cmd

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"inspect/pkg/configs"
	"inspect/pkg/global"
	"os"
	"path"
)

type errorOnWarningHook struct{}

func (errorOnWarningHook) Levels() []log.Level {
	return []log.Level{log.WarnLevel}
}

func (errorOnWarningHook) Fire(entry *log.Entry) error {
	log.Fatalln(entry.Message)
	return nil
}

var (
	GlobalVerbose          bool
	GlobalSuppressWarnings bool
	GlobalErrorOnWarning   bool
	GlobalFiles            []string
)
var RootCmd = &cobra.Command{
	Use:           "cetools",
	Short:         "CloudExplorer 的巡检和高可用部署工具",
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Add extra logging when verbosity is passed
		if GlobalVerbose {
			log.SetLevel(log.DebugLevel)
		}

		// Disable the timestamp (too fast!)
		formatter := new(log.TextFormatter)
		formatter.DisableTimestamp = true
		formatter.ForceColors = true
		log.SetFormatter(formatter)

		if GlobalSuppressWarnings {
			log.SetLevel(log.ErrorLevel)
		} else if GlobalErrorOnWarning {
			hook := errorOnWarningHook{}
			log.AddHook(hook)
		}
		v := viper.New()
		//v.SetConfigType("yml")
		v.SetConfigName("app")

		for _, file := range GlobalFiles {
			v.AddConfigPath(path.Join(file))
		}

		if err := v.ReadInConfig(); err != nil {
			//panic(err)
		}
		if _, err := os.Stat(v.ConfigFileUsed()); err != nil {
			if os.IsNotExist(err) {
				//fmt.Println("配置文件不存在，跳过加载和解析")
				return
			}
		}

		v.OnConfigChange(func(e fsnotify.Event) {
			if err := v.Unmarshal(&global.CONF); err != nil {
				panic(err)
			}
		})
		serverConfig := configs.ServerConfig{}
		if err := v.Unmarshal(&serverConfig); err != nil {
			panic(err)
		}
		global.CONF = &serverConfig
		//cmd.Flags().VisitAll(func(f *pflag.Flag) {
		//	configName := f.Name
		//	if configName == "file" && !f.Changed && v.IsSet(configName) {
		//		GlobalFiles = v.GetStringSlice(configName)
		//		v.OnConfigChange(func(e fsnotify.Event) {
		//			if err := v.Unmarshal(&global.CONF); err != nil {
		//				panic(err)
		//			}
		//		})
		//
		//		//v.OnConfigChange(func(e fsnotify.Event) {
		//		//	if err := v.Unmarshal(&global.CONF); err != nil {
		//		//		panic(err)
		//		//	}
		//		//})
		//		serverConfig := configs.ServerConfig{}
		//		if err := v.Unmarshal(&serverConfig); err != nil {
		//			panic(err)
		//		}
		//		if fileOp.Stat("~/.cetools/config.yaml") {
		//			//if serverConfig.System.BaseDir != "" {
		//			//	baseDir = serverConfig.System.BaseDir
		//			//}
		//			//if serverConfig.System.Port != "" {
		//			//	port = serverConfig.System.Port
		//			//}
		//			//if serverConfig.System.Version != "" {
		//			//	version = serverConfig.System.Version
		//			//}
		//			//if serverConfig.System.Username != "" {
		//			//	username = serverConfig.System.Username
		//			//}
		//			//if serverConfig.System.Password != "" {
		//			//	password = serverConfig.System.Password
		//			//}
		//			//if serverConfig.System.Entrance != "" {
		//			//	entrance = serverConfig.System.Entrance
		//			//}
		//		}
		//
		//	}
		//})
	},
}

// Execute executes the root level command.
// It returns an error if any.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&GlobalVerbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&GlobalSuppressWarnings, "suppress-warnings", false, "Suppress all warnings")
	RootCmd.PersistentFlags().BoolVar(&GlobalErrorOnWarning, "error-on-warning", false, "Treat any warning as an error")
	RootCmd.PersistentFlags().StringSliceVarP(&GlobalFiles, "file", "f", []string{}, "config file")
}
