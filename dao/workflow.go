package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s-platform-api/db"
	"k8s-platform-api/model"
)

var Workflow workflow

type workflow struct {
}

type WorkflowResp struct {
	Items []*model.Workflow
	total int
}

// 获取workflow列表

func (w *workflow) GetWorkflows(filtername, namespace string, limit, page int) (data *WorkflowResp, err error) {

	//	定义分页的起始位置
	startSet := (page - 1) * limit
	//	定义数据库查询返回的内容
	var workflowList []*model.Workflow

	// 数据库查询，Limit方法用于限制条数，offset方法用于设置起始位置
	// workflowList 需要返回，为了保证是同一个值，这个东西需要是指针
	tx := db.GORM.Where("name like ?", "%"+filtername+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)

	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow列表失败," + tx.Error.Error())
		return nil, errors.New("获取workflow列表失败" + tx.Error.Error())
	}
	return &WorkflowResp{

		Items: workflowList,
		total: len(workflowList),
	}, nil

}

// 获取单条数据

func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {

	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow单条数据失败," + tx.Error.Error())
		return nil, errors.New("获取workflow单条数据失败" + tx.Error.Error())
	}
	return workflow, nil

}

// 表数据新增
func (w *workflow) Add(workflow *model.Workflow) (err error) {

	tx := db.GORM.Create(&workflow)
	if tx.Error != nil {
		logger.Error("添加workflow失败" + tx.Error.Error())
		return errors.New("添加workflow失败" + tx.Error.Error())
	}
	return nil

}

// 表数据删除

func (w *workflow) DelById(id int) (err error) {

	// 这里可以了解一下软硬删除  Delete 是软删除    unscope是硬删除
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil {
		logger.Error("删除workflow失败" + tx.Error.Error())
		return errors.New("删除workflow失败" + tx.Error.Error())
	}

	return nil
}
