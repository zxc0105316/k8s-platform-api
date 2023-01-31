package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}

var Service service

type service struct {
}

func (s *service) CreateService(data *ServiceCreate) (err error) {
	/*
		下面这个配置对标的yaml文件内容

		apiServer: appsv1
		kind: service
		metadata:
			name:
			namespace:
		spec:
			selector:
				app:
			ports:
			- name:
		      port:
			  targerPort:
		      NodePort:


	*/

	//	将data中的数据组装成一个corev1.service对象
	service := &corev1.Service{
		//TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
		Status: corev1.ServiceStatus{},
	}
	// 默认是clusterIP，这里做外部判断Nodeport，添加配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}
	//	 创建service
	_, err = K8s.ClientSet.CoreV1().Services(data.Name).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {

		logger.Error(errors.New("创建service失败," + err.Error()))
		return errors.New("创建service失败" + err.Error())

	}

	return nil

}

func (s *service) DelService(ServiceName, Namespace string) (err error) {

	err = K8s.ClientSet.NetworkingV1().Ingresses(Namespace).Delete(context.TODO(), ServiceName, metav1.DeleteOptions{})

	if err != nil {

		logger.Error("删除ingress失败" + err.Error())
		return errors.New("删除ingress失败" + err.Error())

	}
	return nil
}

// 类型转换的方法   corev1.Service -> DataCell , DataCell ->  corev1.Service

func (s *service) toCells(services []corev1.Service) []DataCell {

	cells := make([]DataCell, len(services))

	for i := range services {
		cells[i] = serviceCell(services[i])
	}
	return cells
}

// DataCell -> corev1.Service

func (s *service) fromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		// 这里用断言反转换
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}

// 定义列表的返回内容,total是元素数量,Items是service元素列表
type ServicesResp struct {
	Total int              `json:"total"`
	Items []corev1.Service `json:"items"`
}

// 获取service列表，支持过滤，排序和分页

func (s *service) GetServices(filterName, namespace string, limit, page int) (data *ServicesResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	ServicesList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Error("获取Service列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取Service列表失败," + err.Error())
	}

	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: s.toCells(ServicesList.Items),
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

	services := s.fromCells(dataList.GenericDatalist)

	return &ServicesResp{
		Total: total,
		Items: services,
	}, nil

}

// 获取service详情
func (s *service) GetServicesDetail(serviceName, namespace string) (service *corev1.Service, err error) {
	thisservice, err := K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取service详情失败," + err.Error())
		return nil, errors.New("获取service失败," + err.Error())
	}

	return thisservice, err
}

// 删除service
func (s *service) DeleteService(serviceName, namespace string) (err error) {

	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})

	if err != nil {
		logger.Error("删除service失败," + err.Error())
		return errors.New("删除service失败," + err.Error())
	}

	return nil

}

// 更新service
// context 参数是请求中传入的service对象的json数组,是对象的序列化  marshal
// 注意这个service对象是corev1.service类型的
func (s *service) UpdateService(serviceName, namespace, Content string) error {

	var service = &corev1.Service{}

	//反序列化
	err := json.Unmarshal([]byte(Content), service)
	if err != nil {
		logger.Error(errors.New("反序列化失败," + err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}

	//	更新service
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {

		logger.Error(errors.New("更新service失败," + err.Error()))
		return errors.New("更新service失败," + err.Error())

	}
	return nil

}

// 获取每个namespace下的service数量

type ServicesNp struct {
	Namespace  string
	ServiceNum int
}

func (p *service) GetServiceNumberNp() (servicesNps []*ServicesNp, err error) {

	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		//	获取service列表
		serviceList, err := K8s.ClientSet.CoreV1().Services(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//	 组装数据
		servicesNp := &ServicesNp{
			Namespace:  namespace.Name,
			ServiceNum: len(serviceList.Items),
		}
		//	 添加到servicesNps数组中
		servicesNps = append(servicesNps, servicesNp)
	}

	return servicesNps, nil
}
