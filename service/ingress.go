package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress ingress

type ingress struct {
}

/*对标yaml

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
	name:
	namespace:
spec:
	rules:
	- host: www.sss.com
      http:
        paths: /
		- pathType: prefix
		  backend:
			Ingress:
 				name: myapp-svc
				port:
					number: 80

*/

type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace""`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

// 定义ingress的path结构体
type HttpPath struct {
	Path          string        `json:"path"`
	PathType      nwv1.PathType `json:"path_type"`
	IngressesName string        `json:"Ingress_name"`
	IngressesPort int32         `json:"Ingress_porti"`
}

// 创建ingress
func (i *ingress) CreateIngress(data *IngressCreate) (err error) {

	//	声明nwv1.IngressRule和nwv1.HTTPIngressPath对象
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath
	//	将data中的数据组装成nwv1.Ingress对象

	ingress := &nwv1.Ingress{
		//TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Status: nwv1.IngressStatus{},
	}

	// 套两层循环
	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{
					Paths: nil,
				},
			},
		}
		for _, HttpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     HttpPath.Path,
				PathType: &HttpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: HttpPath.IngressesName,
						Port: nwv1.ServiceBackendPort{
							Number: HttpPath.IngressesPort,
						},
					},
					Resource: nil,
				},
			}
			//	 封装hip为数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		ingressRules = append(ingressRules, ir)

	}
	ingress.Spec.Rules = ingressRules

	// 创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建ingress失败, " + err.Error()))
		errors.New("创建ingress失败" + err.Error())
	}
	return nil
}

// 获取Ingress详情

func (i *ingress) GetIngressessDetail(IngressName, namespace string) (Ingress *nwv1.Ingress, err error) {
	thisIngress, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), IngressName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Ingress详情失败," + err.Error())
		return nil, errors.New("获取Ingress失败," + err.Error())
	}

	return thisIngress, err
}

// 删除ingress

func (i *ingress) DeleteIngresses(IngressName, namespace string) (err error) {

	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), IngressName, metav1.DeleteOptions{})

	if err != nil {
		logger.Error("删除Ingress失败," + err.Error())
		return errors.New("删除Ingress失败," + err.Error())
	}

	return nil

}

// 更新ingress

func (i *ingress) UpdateIngresses(IngressName, namespace, Content string) error {

	var Ingress = &nwv1.Ingress{}

	//反序列化
	err := json.Unmarshal([]byte(Content), Ingress)
	if err != nil {
		logger.Error(errors.New("反序列化失败," + err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}

	//	更新Ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), Ingress, metav1.UpdateOptions{})
	if err != nil {

		logger.Error(errors.New("更新Ingress失败," + err.Error()))
		return errors.New("更新Ingress失败," + err.Error())

	}
	return nil

}

// 定义列表的返回内容,total是元素数量,Items是service元素列表
type IngressesResp struct {
	Total int            `json:"total"`
	Items []nwv1.Ingress `json:"items"`
}

// 获取service列表，支持过滤，排序和分页

func (i *ingress) GetIngresses(filterName, namespace string, limit, page int) (data *IngressesResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	IngressesList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Error("获取Ingresses列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取Ingresses列表失败," + err.Error())
	}

	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: i.toCells(IngressesList.Items),
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

	Ingresses := i.fromCells(dataList.GenericDatalist)

	return &IngressesResp{
		Total: total,
		Items: Ingresses,
	}, nil

}

// 类型转换的方法   nwv1.ingress -> DataCell , DataCell ->  nwv1.ingress

func (i *ingress) toCells(Ingresses []nwv1.Ingress) []DataCell {

	cells := make([]DataCell, len(Ingresses))

	for i := range Ingresses {
		cells[i] = ingressCell(Ingresses[i])
	}
	return cells
}

// DataCell -> nwv1.ingress

func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresses := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		// 这里用断言反转换
		ingresses[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingresses
}
