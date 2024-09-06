package viper

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"inspect/pkg/configs"
	"inspect/pkg/global"
	"os"
	"path"
)

func Init() {

	v := viper.New()
	//v.SetConfigType("yml")
	v.SetConfigName("app")

	for _, file := range global.Conf.Files {
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

}
