package global

import (
	"cetool/pkg/configs"
	"cetool/pkg/dto"
	"cetool/pkg/utils/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
