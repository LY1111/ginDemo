package controllers

import (
    "github.com/gin-gonic/gin"
    "tag_data_sync/models/http"
)

/**
 * @Author: liyan
 * @Description: TODO
 * @File:  no_router
 * @Version: 1.0.0
 * @Date: 2021/5/18 12:14 下午
 */
func NoRouter(c *gin.Context) {
    http.FailWithMessage(c,"请求地址错误")
}