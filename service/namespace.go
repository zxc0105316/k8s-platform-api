package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*原始yaml

apiVersion: v1
kind: Namespace
metadata:
	name: test1
	label:
		name: test1
spec:


*/

type Namespace namespace

type namespace struct {
}

// 创建namespace
func (n *namespace) CreateNamespace(namespace string) (err error) {

	Namespace := &corev1.Namespace{
		//TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}

	// 创建namespace
	_, err = K8s.ClientSet.CoreV1().Namespaces().Create(context.TODO(), Namespace, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Namespace失败, " + err.Error()))
		errors.New("创建Namespace失败" + err.Error())
	}
	return nil
}

// 获取namespace详情

func (n *namespace) GetNamespaceDetail(namespace string) (Namespace *corev1.Namespace, err error) {
	thisnamespace, err := K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取namespace详情失败," + err.Error())
		return nil, errors.New("获取namespace失败," + err.Error())
	}

	return thisnamespace, err
}

// 删除namespace

func (n *namespace) DeleteNamespaces(namespace string) (err error) {

	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})

	if err != nil {
		logger.Error("删除namespace失败," + err.Error())
		return errors.New("删除namespace失败," + err.Error())
	}

	return nil

}

// 更新ingress

func (i *ingress) UpdateNamespace(namespace, Content string) error {

	var Namespace = &corev1.Namespace{}

	//反序列化
	err := json.Unmarshal([]byte(Content), Namespace)
	if err != nil {
		logger.Error(errors.New("反序列化失败," + err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}

	//	更新Namespace
	_, err = K8s.ClientSet.CoreV1().Namespaces().Update(context.TODO(), Namespace, metav1.UpdateOptions{})
	if err != nil {

		logger.Error(errors.New("更新Namespace失败," + err.Error()))
		return errors.New("更新Namespace失败," + err.Error())

	}
	return nil

}

// 定义列表的返回内容,total是元素数量,Items是namespace元素列表
type NamespacesResp struct {
	Total int                `json:"total"`
	Items []corev1.Namespace `json:"items"`
}

// 获取namespace列表，支持过滤，排序和分页

func (n *namespace) GetNamespaces(limit, page int) (data *NamespacesResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	NamespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Error("获取Namespace列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取Namespace列表失败," + err.Error())
	}
	// 命名空间的查询不需要过滤条件
	filterName := ""
	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: n.toCells(NamespaceList.Items),
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

	Namespaces := n.fromCells(dataList.GenericDatalist)

	return &NamespacesResp{
		Total: total,
		Items: Namespaces,
	}, nil

}

// 类型转换的方法   corev1.Namespace -> DataCell , DataCell ->  corev1.Namespace

func (n *namespace) toCells(Namespaces []corev1.Namespace) []DataCell {

	cells := make([]DataCell, len(Namespaces))

	for i := range Namespaces {
		cells[i] = NamespaceCell(Namespaces[i])
	}
	return cells
}

// DataCell -> corev1.Namespace

func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	Namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		// 这里用断言反转换
		Namespaces[i] = corev1.Namespace(cells[i].(NamespaceCell))
	}
	return Namespaces
}
