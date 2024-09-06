package global

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"inspect/pkg/configs"
	"inspect/pkg/utils/logger"
)

var (
	LOG   *logrus.Logger
	Print *logger.Logger
	CONF  *configs.ServerConfig
	Viper *viper.Viper
)
