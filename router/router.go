package router

import (
    "github.com/gin-gonic/gin"
    "tag_data_sync/controllers"
)

func LoadRouter() *gin.Engine {
    r := gin.New()

    r.NoRoute(controllers.NoRouter)



    return r
}
