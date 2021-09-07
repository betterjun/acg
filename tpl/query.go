package tpl

import (
	"os"

	"github.com/spf13/viper"
)

func ParseConfigure(file string) (*PackageCfg, error) {
	config := viper.New()
	if len(file) == 0 {
		//获取项目的执行路径
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		config.AddConfigPath(path)     //设置读取的文件路径
		config.SetConfigName("cr_all") //设置读取的文件名
		config.SetConfigType("yaml")   //设置文件的类型
	} else {
		config.SetConfigFile(file)
	}

	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	pkgCfg := NewPackageCfg()
	err := pkgCfg.parsePackage(config)
	if err != nil {
		panic(err)
	}

	return pkgCfg, nil
}
