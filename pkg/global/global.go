package global

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"inspect/pkg/configs"
	"inspect/pkg/dto"
	"inspect/pkg/utils/logger"
)

var (
	Conf dto.GlobalConf
)

var (
	LOG   *logrus.Logger
	Print *logger.Logger
	CONF  *configs.ServerConfig
	Viper *viper.Viper
)
