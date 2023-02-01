package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PersistentVolume persistentVolume

type persistentVolume struct {
}

type PersistentVolumeCreate struct {
	Name                          string                       `json:"name"`
	storageClassName              string                       `json:"storageClassName"`
	accessModes                   []string                     `json:"accessModes"`
	VolumeMode                    *corev1.PersistentVolumeMode `json:"volume_mode"`
	capacity                      map[string]string            `json:"capacity"`
	PersistentVolumeReclaimPolicy string                       `json:"persistent_volume_reclaim_policy"`
	NodeAffinity                  *corev1.VolumeNodeAffinity   `json:"node_affinity"`
	MountOptions                  []string                     `json:"mountOptions"`
	NodeSelectorTerm              []string
}

// 创建pv

//func (p *persistentVolume) CreatePersistentVolume(data *PersistentVolumeCreate) (err error) {
//
//	pv := corev1.PersistentVolume{
//		//TypeMeta:   metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name: data.Name,
//		},
//		Spec: corev1.PersistentVolumeSpec{
//			Capacity:                      corev1.ResourceList{},
//			PersistentVolumeSource:        corev1.PersistentVolumeSource{},
//			AccessModes:                   nil,
//			ClaimRef:                      nil,
//			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimPolicy(data.PersistentVolumeReclaimPolicy),
//			StorageClassName:              data.storageClassName,
//			MountOptions:                  nil, // []string
//			VolumeMode:                    &corev1.PersistentVolumeMode(),
//			NodeAffinity: &corev1.VolumeNodeAffinity{
//				Required: &corev1.NodeSelector{
//					NodeSelectorTerms: []corev1.NodeSelectorTerm{
//						MatchExpressions: data.NodeSelectorTerm[MatchExpressions],
//						MatchFields:      nil,
//					},
//				}},
//		},
//	}
//
//}

// 修改pv,更新pv

// 获取pv

// 删除pv

// 获取pvlist

// 定义列表的返回内容,total是元素数量,Items是namespace元素列表
type PersistentVolumesResp struct {
	Total int                       `json:"total"`
	Items []corev1.PersistentVolume `json:"items"`
}

// 获取PersistentVolume列表，支持过滤，排序和分页

func (p *persistentVolume) GetPersistentVolumes(limit, page int) (data *PersistentVolumesResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	PersistentVolumeList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Error("获取PersistentVolume列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取PersistentVolume列表失败," + err.Error())
	}
	// pv的查询不需要过滤条件
	filterName := ""
	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: p.toCells(PersistentVolumeList.Items),
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

	PersistentVolumes := p.fromCells(dataList.GenericDatalist)

	return &PersistentVolumesResp{
		Total: total,
		Items: PersistentVolumes,
	}, nil

}

// 类型转换的方法   corev1.PersistentVolume -> DataCell , DataCell ->  corev1.PersistentVolume

func (p *persistentVolume) toCells(PersistentVolumes []corev1.PersistentVolume) []DataCell {

	cells := make([]DataCell, len(PersistentVolumes))

	for i := range PersistentVolumes {
		cells[i] = persistentVolumeCell(PersistentVolumes[i])
	}
	return cells
}

// DataCell -> corev1.PersistentVolume

func (p *persistentVolume) fromCells(cells []DataCell) []corev1.PersistentVolume {
	PersistentVolumes := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		// 这里用断言反转换
		PersistentVolumes[i] = corev1.PersistentVolume(cells[i].(persistentVolumeCell))
	}
	return PersistentVolumes
}
