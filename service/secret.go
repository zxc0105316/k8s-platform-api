package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct {
}

type SecretCreate struct {
}

// 创建

func (s *secret) CreateSecret(data *SecretCreate) (err error) {

	return nil

}

// 删除
func (s *secret) DelSecret(secretName, namespace string) (err error) {

	return nil
}

// 更新
func (s *secret) UpdateSecret(secretName, namespace, content string) (err error) {

	return nil

}

// 获取详情

func (s *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {

	tsecret := &corev1.Secret{}

	return tsecret, nil
}

// 定义列表的返回内容,total是元素数量,Items是Secret元素列表
type SecretsResp struct {
	Total int             `json:"total"`
	Items []corev1.Secret `json:"items"`
}

// 获取namespace列表，支持过滤，排序和分页

func (s *secret) GetSecrets(namespace string, limit, page int) (data *SecretsResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	SecretList, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
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

		GenericDatalist: s.toCells(SecretList.Items),
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

	Secrets := s.fromCells(dataList.GenericDatalist)

	return &SecretsResp{
		Total: total,
		Items: Secrets,
	}, nil

}

// 类型转换的方法   corev1.Secret -> DataCell , DataCell ->  corev1.Secret

func (s *secret) toCells(Secrets []corev1.Secret) []DataCell {

	cells := make([]DataCell, len(Secrets))

	for i := range Secrets {
		cells[i] = secretCell(Secrets[i])
	}
	return cells
}

// DataCell -> corev1.Namespace

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	Secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		// 这里用断言反转换
		Secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return Secrets
}
