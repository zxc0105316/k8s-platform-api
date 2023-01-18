package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct {
}

// 类型转换的方法   corev1.Pod -> DataCell , DataCell ->  corev1.Pod

func (p *pod) toCells(pods []corev1.Pod) []DataCell {

	cells := make([]DataCell, len(pods))

	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// DataCell -> corev1.Pod

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		// 这里用断言反转换
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// 定义列表的返回内容,total是元素数量,Items是pod元素列表
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

// 获取pod列表，支持过滤，排序和分页

func (p *pod) GetPods(filterName, namespace string, limit, page int) (data *PodsResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	podsList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Info("获取Pod列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取Pod列表失败," + err.Error())
	}

	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: p.toCells(podsList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{
				Name: filterName,
			},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDatalist)

	// 排序和分页
	dataList := filtered.Sort().Paginate()

	//将DataCell类型转换

	pods := p.fromCells(dataList.GenericDatalist)

	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil

}
