package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform-api/config"
	"k8s-platform-api/controller"
	"k8s-platform-api/db"
	"k8s-platform-api/service"
)

func main() {

	// 初始化数据库
	db.Init()

	// 初始化 k8s client

	service.K8s.Init()

	// 初始化gin
	r := gin.Default()

	controller.Router.InitApiRouter(r)

	r.Run(config.ListenAddr)

	//	关闭数据库连接
	db.Close()

}
