package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var Deployment deployment

type deployment struct {
}

// 创建deployment
func (d *deployment) CreateDeployment(ctx *gin.Context) {

	// 组装deployCreate 数据
	data := &service.DeployCreate{}

	if err := ctx.ShouldBind(data); err != nil {
		logger.Error(errors.New("创建deployment: 绑定数据失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "创建deployment: 绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Deployment.CreateDeployment(data)
	if err != nil {

		logger.Error(errors.New("创建deployment失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建deployment失败," + err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建deployment成功",
		"data": nil,
	})

}

// 删除deployment
func (d *deployment) DelDeployment(ctx *gin.Context) {

	//组装数据
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
	})
	if err := ctx.ShouldBind(params); err != nil {

		logger.Error(errors.New("删除Deployment: 绑定数据失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "删除deployment失败：" + err.Error(),
			"data": nil,
		})
		return

	}

	err := service.Deployment.DelDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		logger.Error("删除deployment失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除deployment失败" + err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除deployment成功",
		"data": nil,
	})
}

// 重启deployment
func (d *deployment) RestartDeployment(ctx *gin.Context) {

	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		logger.Error("重启deployment： 绑定数据失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "重启deployment：绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)

	if err != nil {
		logger.Error("重启deployment失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "重启deployment失败" + err.Error(),
			"data": nil,
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "重启deployment成功",
		"data": nil,
	})

}

// 修改deployment副本数
func (d *deployment) UpdateDeploymentReplicas(ctx *gin.Context) {

	//组装数据
	params := new(struct {
		DeploymentName string `json:"deploymentName"`
		Namespace      string `json:"namespace"`
		Replicas       int    `json:"replicas"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		logger.Error("修改deployment副本数： 绑定数据失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "修改deployment副本数：绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.Replicas)

	if err != nil {
		logger.Error("修改deployment副本数失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "重启deployment失败" + err.Error(),
			"data": nil,
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "修改deployment副本数成功",
		"data": fmt.Sprintf("最新副本数：%d", data),
	})

}

// 更新deployment
func (d *deployment) UpdateDeployment(ctx *gin.Context) {

	params := new(struct {
		podName   string `json:"podName"`
		namespace string `json:"namespace"`
		content   string `json:"content"`
	})

	if err := ctx.ShouldBind(params); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "更新deployment: 绑定数据失败",
			"data": nil,
		})
		return

	}

	err := service.Deployment.UpdateDeployment(params.podName, params.namespace, params.content)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新deployment成功",
		"data": nil,
	})

}

// 获取deployment列表，支持分页，过滤，排序
func (d *deployment) GetDeployments(ctx *gin.Context) {

	//	处理入参
	//  匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"deploymentName"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用bind方法，json格式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("获取deployment: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取deployment: Bind绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeployments(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment列表成功",
		"data": data,
	})
}

// 获取deployment详情

func (d *deployment) GetDeploymentDetail(ctx *gin.Context) {

	//	处理入参
	//	构造匿名结构体,用于判断传入的参数
	params := new(struct {
		PodName   string `form:"DeploymentName"`
		Namespace string `form:"namespace"`
	})

	// form 方法使用Bind()方法，json方式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {

		logger.Error("获取指定deployment详情: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取指定deployment详情: Bind绑定失败",
			"data": nil,
		})
		return
	}
	deployment, err := service.Deployment.GetDeploymentsDetail(params.PodName, params.Namespace)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{

		"msg":  "获取pod详情成功",
		"data": deployment,
	})
}

// 获取每个namespace下的pod数量

func (d *deployment) GetDeploymentNumberNp(ctx *gin.Context) {

	// 无参
	//params := new(struct{})
	data, err := service.Deployment.GetDeploymentNumberNp()

	fmt.Println(data)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取namespace->deployment->count 失败",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个namespace下deployment数量成功",
		"data": data,
	})

}
