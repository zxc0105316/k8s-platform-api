package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

// 处理ctx

// 全局导出pod
var Pod pod

type pod struct {
}

// 获取pod列表，支持分页，过滤，排序
func (p *pod) GetPods(ctx *gin.Context) {

	//	处理入参
	//  匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"podName"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用bind方法，json格式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("获取pod: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取pod: Bind绑定失败" + err.Error(),
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

// 获取pod详情

func (p *pod) GetPodDetail(ctx *gin.Context) {

	//	处理入参
	//	构造匿名结构体,用于判断传入的参数
	params := new(struct {
		PodName   string `form:"podName"`
		Namespace string `form:"namespace"`
	})

	// form 方法使用Bind()方法，json方式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {

		logger.Error("获取指定pod详情: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取指定pod详情: Bind绑定失败",
			"data": nil,
		})
		return
	}
	pod, err := service.Pod.GetPodsDetail(params.PodName, params.Namespace)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{

		"msg":  "获取pod详情成功",
		"data": pod,
	})
}

// 删除指定pod

func (p *pod) DeletePod(ctx *gin.Context) {

	//	处理入参
	//	构造匿名结构体,用于判断传入的参数
	params := new(struct {
		PodName   string `json:"podName"`
		Namespace string `json:"namespace"`
	})

	// form 方法使用Bind()方法，json方式使用shouldBindJSON方法
	if err := ctx.ShouldBind(params); err != nil {

		logger.Error("删除pod: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除pod: Bind绑定失败",
			"data": nil,
		})
		return
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{

		"msg":  "删除pod成功",
		"data": nil,
	})
}

// 更新指定pod

func (p *pod) UpdatePod(ctx *gin.Context) {

	params := new(struct {
		podName   string `json:"podName"`
		namespace string `json:"namespace"`
		content   string `json:"content"`
	})

	if err := ctx.ShouldBind(params); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "更新pod: 绑定数据失败",
			"data": nil,
		})
		return

	}

	err := service.Pod.UpdatePod(params.podName, params.namespace, params.content)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})
}

// 获取pod中的容器名字

func (p *pod) GetPodContainer(ctx *gin.Context) {

	params := new(struct {
		podName   string `form:"podName"`
		namespace string `form:"namespace"`
	})

	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "获取pod中容器名: 绑定失败",
			"data": nil,
		})
		return
	}
	containersList, err := service.Pod.GetPodContainer(params.podName, params.namespace)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "获取pod内容器列表失败 ," + err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod内容器列表成功",
		"data": containersList,
	})

}

// 获取pod中的日志
func (p *pod) GetPodLog(ctx *gin.Context) {

	params := new(struct {
		containerName string `form:"containerName"`
		podName       string `form:"podName"`
		namespace     string `form:"namespace"`
	})

	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "获取pod日志: 绑定失败",
			"data": nil,
		})
		return
	}
	log, err := service.Pod.GetPodLog(params.podName, params.namespace, params.containerName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "获取pod" + params.containerName + " 日志失败+," + err.Error(),
			"data": log,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod内容器名成功",
		"data": nil,
	})

}

// 获取每个namespace下的pod数量

func (p *pod) GetPodNumberNp(ctx *gin.Context) {

	// 无参
	//params := new(struct{})
	data, err := service.Pod.GetPodNumberNp()

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取namespace->pod->count 失败",
			"data": data,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取namespace下pod数量成功",
		"data": data,
	})

}
