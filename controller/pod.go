package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform-api/service"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

// 处理ctx

var Pod v1.Pod

type pod struct {
}

// 获取pod列表，支持分页，过滤，排序
func (p *pod) GetPods(ctx *gin.Context) {

	//	处理入参
	//  匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用bind方法，json格式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}
