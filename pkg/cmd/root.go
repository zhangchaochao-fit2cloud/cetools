package cmd

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"inspect/pkg/global"
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
		v.SetConfigType("yaml")
		_ = v.BindEnv("file", "CONFIG_FILE")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			configName := f.Name
			if configName == "file" && !f.Changed && v.IsSet(configName) {
				GlobalFiles = v.GetStringSlice(configName)
				v.OnConfigChange(func(e fsnotify.Event) {
					if err := v.Unmarshal(&global.CONF); err != nil {
						panic(err)
					}
				})

			}
		})
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
}
