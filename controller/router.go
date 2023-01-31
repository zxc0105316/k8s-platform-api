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

	router.
		GET("/test", func(context *gin.Context) {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "测试",
				"data": "测试",
			})
		}).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		POST("/api/k8s/pod/delete", Pod.DeletePod).
		POST("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		GET("/api/k8s/pod/ContainerList", Pod.GetPodContainer).
		GET("/api/k8s/pod/NumberNp", Pod.GetPodNumberNp).
		GET("/api/k8s/pod/podLists", Pod.GetPods).
		POST("/api/k8s/deployment/delete", Deployment.DelDeployment).
		POST("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		GET("/api/k8s/deployment/NumberNp", Deployment.GetDeploymentNumberNp).
		GET("/api/k8s/deployment/deploymentLists", Deployment.GetDeployments)
}
