package global

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"inspect/pkg/common"
	"inspect/pkg/configs"
)

var (
	LOG   *logrus.Logger
	Print *common.Logger
	CONF  configs.ServerConfig
	Viper *viper.Viper
)
