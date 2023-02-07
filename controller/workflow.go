package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform-api/service"
	"net/http"
)

var Workflow workflow

type workflow struct {
}

// 获取分页数据

func (w *workflow) GetWorkflowList(ctx *gin.Context) {

	//	处理入参
	//  匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用bind方法，json格式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("获取workflow: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取workflow: Bind绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Workflow.GetList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取workflow列表成功",
		"data": data,
	})
}

// 获取单条数据

func (w *workflow) GetWorkflow(ctx *gin.Context) {

	//	处理入参
	//	构造匿名结构体,用于判断传入的参数
	params := new(struct {
		id int `form:"id"`
	})

	// form 方法使用Bind()方法，json方式使用shouldBindJSON方法
	if err := ctx.Bind(params); err != nil {

		logger.Error("获取指定workflow详情: Bind绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取指定workflow详情: Bind绑定失败",
			"data": nil,
		})
		return
	}
	workflow, err := service.Workflow.GetWorkflow(params.id)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{

		"msg":  "获取workflow详情成功",
		"data": workflow,
	})

}

// 创建workflow

func (w *workflow) CreateWorkflow(ctx *gin.Context) {

	// 组装workflowCreate 数据
	data := &service.WorkflowCreate{}

	if err := ctx.ShouldBind(data); err != nil {
		logger.Error(errors.New("创建workflow: 绑定数据失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "创建workflow: 绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Workflow.CreateWorkflow(data)
	if err != nil {

		logger.Error(errors.New("创建workflow失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建workflow失败," + err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建workflow成功",
		"data": nil,
	})

}

// 删除workflow

func (w *workflow) DelWorkflow(ctx *gin.Context) {

	//组装数据
	params := new(struct {
		id int `json:"id"`
	})
	if err := ctx.ShouldBind(params); err != nil {

		logger.Error(errors.New("删除Workflow: 绑定数据失败"))
		ctx.JSON(http.StatusInternalServerError, gin.H{

			"msg":  "删除Workflow失败：" + err.Error(),
			"data": nil,
		})
		return

	}

	err := service.Workflow.DelById(params.id)
	if err != nil {
		logger.Error("删除Workflow失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除Workflow失败" + err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Workflow成功",
		"data": nil,
	})
}
