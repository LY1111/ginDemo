package main

import (
    "context"
    "data_binding_backend/config"
    "data_binding_backend/initialize/logger"
    "data_binding_backend/initialize/mysql"
    _ "data_binding_backend/initialize/viper"
    "data_binding_backend/router"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func init() {
    // 初始化log日志
    _ = logger.InitLogger(&config.Config.Log)
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
            logger.Errorf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit

    _ = mysql.CloseMysqlPool()

    logger.Error("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server Shutdown:", err)
    }
    logger.Error("Server exiting")
}
