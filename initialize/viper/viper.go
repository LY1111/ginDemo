package viper

import (
    "data_binding_backend/config"
    "data_binding_backend/initialize/logger"
    "flag"
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

func init() {
    fmt.Println("viper")
    var err error
    configPath := flag.String("f", "./config/config_dev.toml", "config file path error")
    flag.Parse()

    viper.SetConfigType("toml")
    viper.SetConfigFile(*configPath)
    err = viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("parse config err:%s", err.Error()))
    }

    _ = viper.Unmarshal(&config.Config)
    viper.WatchConfig()
    viper.OnConfigChange(func(in fsnotify.Event) {
        logger.Info("Config file change: ", in.Name)
        _ = viper.Unmarshal(&config.Config)
    })
}
