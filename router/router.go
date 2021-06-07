package router

import (
    "data_binding_backend/controllers"
    "github.com/gin-gonic/gin"
)

func LoadRouter() *gin.Engine {
    r := gin.New()

    r.NoRoute(controllers.NoRouter)



    return r
}
