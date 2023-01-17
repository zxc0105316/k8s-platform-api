package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义一个全局的路由变量，用于跨包调用
var Router router

// 定义一个router结构体
type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {

	router.GET("/testapi", func(context *gin.Context) {

		context.JSON(http.StatusOK, gin.H{
			"msg":  "test success!",
			"data": nil,
		})
	})

}
