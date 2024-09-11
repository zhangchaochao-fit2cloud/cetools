package log

import (
	"cetool/pkg/global"
	"github.com/sirupsen/logrus"
)

type errorOnWarningHook struct{}

func (errorOnWarningHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel}
}

func (errorOnWarningHook) Fire(entry *logrus.Entry) error {
	logrus.Fatalln(entry.Message)
	return nil
}

func Init() {
	logger := logrus.New()
	global.LOG = logger

	if global.Conf.Verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Disable the timestamp (too fast!)
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.ForceColors = true
	logger.SetFormatter(formatter)

	if global.Conf.SuppressWarnings {
		logger.SetLevel(logrus.ErrorLevel)
	} else if global.Conf.ErrorOnWarning {
		hook := errorOnWarningHook{}
		logger.AddHook(hook)
	}
}

//type MineFormatter struct{}
//func (s *MineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
//	detailInfo := ""
//	if entry.Caller != nil {
//		function := strings.ReplaceAll(entry.Caller.Function, "github.com/1Panel-dev/1Panel/backend/", "")
//		detailInfo = fmt.Sprintf("(%s: %d)", function, entry.Caller.Line)
//	}
//	if len(entry.Data) == 0 {
//		msg := fmt.Sprintf("[%s] [%s] %s %s \n", time.Now().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message, detailInfo)
//		return []byte(msg), nil
//	}
//	msg := fmt.Sprintf("[%s] [%s] %s %s {%v} \n", time.Now().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message, detailInfo, entry.Data)
//	return []byte(msg), nil
//}
//
//func setOutput(logger *logrus.Logger, config configs.LogConfig) {
//	writer, err := logrus.NewWriterFromConfig(&log.Config{
//		LogPath:            global.CONF.System.LogPath,
//		FileName:           config.LogName,
//		TimeTagFormat:      FileTImeFormat,
//		MaxRemain:          config.MaxBackup,
//		RollingTimePattern: RollingTimePattern,
//		LogSuffix:          config.LogSuffix,
//	})
//	if err != nil {
//		panic(err)
//	}
//	level, err := logrus.ParseLevel(config.Level)
//	if err != nil {
//		panic(err)
//	}
//	fileAndStdoutWriter := io.MultiWriter(writer, os.Stdout)
//
//	logger.SetOutput(fileAndStdoutWriter)
//	logger.SetLevel(level)
//	logger.SetFormatter(new(MineFormatter))
//}
