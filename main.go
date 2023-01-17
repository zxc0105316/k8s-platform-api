package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform-api/config"
	"k8s-platform-api/controller"
)

func main() {

	r := gin.Default()

	controller.Router.InitApiRouter(r)

	r.Run(config.ListenAddr)

}
