package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "os/signal"
    "tag_data_sync/config"
    "tag_data_sync/initialize/logger"
    "tag_data_sync/initialize/mysql"
    _ "tag_data_sync/initialize/viper"
    "tag_data_sync/router"
    "time"
)

func init() {
    // 初始化log日志
    logger.InitLogger(&config.Config.Log)
    logger.Debugf("log config:%+v", config.Config.Log)
    defer logger.Log.Sync()

    // 初始化MySQL
    err := mysql.InitMySQL()
    if err != nil {
        panic(err)
    }
}

func main() {
    gin.SetMode(config.Config.Server.RunMode)
    srv := &http.Server{
        Addr:              config.Config.Server.HttpPort,
        Handler:           router.LoadRouter(),
        ReadHeaderTimeout: 10 * time.Second,
        WriteTimeout:      10 * time.Second,
        MaxHeaderBytes:    1 << 20,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit

    mysql.CloseMysqlPool()

    logger.Error("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server Shutdown:", err)
    }
    logger.Error("Server exiting")
}
