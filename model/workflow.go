package model

import "time"

// 工作流数据    时间->控制器->service->ingress

type Workflow struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	Type       string `json:"type" gorm:"colume:type"`
	//  Type: clusterip nodeport ingress

}

// 定义tableName方法，返回名mysql表名，以此来自定义mysql中的表名

func (*Workflow) TableName() string {
	return "workflow"
}
