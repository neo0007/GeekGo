package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	//initViperRemote()
	intiViperV1()
	initLogger()
	r := InitWebServer()
	//r.Run(":8081")
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	//如果你不 replace，直接用zap.L() 你啥都打不出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("搞好了")

	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数", zap.Error(errors.New("This is a error")),
		zap.Int64("id", 123),
		zap.Any("结构体", Demo{Name: "hello"}))
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	//通过 webook 和其他使用 etcd 的程序区别出来
	err := viper.AddRemoteProvider("etcd3",
		"127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func intiViperV1() {
	// pflag.String 中 value 是默认值，可以通过环境变量覆盖：
	// go run . --config=config/dev.yaml
	// 或在 GolandIDE 中设置program arguments 参数覆盖：--config=config/dev.yaml
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	//viper.SetConfigFile("config/dev.yaml")
	//实时监听配置变更
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name, e.Op)
		fmt.Println(viper.GetString("db.mysql.dsn"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func intiViper() {
	// 请注意：配置文件的名字不包括文件的扩展名！
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中有很多格式：JSON，XML，YAML，TOML 等
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	//可以有多个 path，可以 逐个扫描读取
	viper.AddConfigPath("./config")
	// 读取配置到 viper，你可以理解就是加载到内存里
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 可以有多个 viper 的实例
	//otherViper := viper.New()
	//otherViper.SetConfigName("myjosn")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}
