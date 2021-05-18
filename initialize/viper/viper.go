package viper

import (
    "flag"
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
    "tag_data_sync/config"
    "tag_data_sync/initialize/logger"
)

func init() {
    var err error
    configPath := flag.String("f", "./config/config_dev.toml", "config file path error")
    flag.Parse()

    viper.SetConfigType("toml")
    viper.SetConfigFile(*configPath)
    err = viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("parse config err:%s", err.Error()))
    }

    viper.Unmarshal(&config.Config)
    viper.WatchConfig()
    viper.OnConfigChange(func(in fsnotify.Event) {
        logger.Info("Config file change: ", in.Name)
        viper.Unmarshal(&config.Config)
    })
}
