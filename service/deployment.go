package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct {
}
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_pathe"`
}

// 定义列表的返回内容,total是元素数量,Items是pod元素列表
type DeploymentsResp struct {
	Total int                 `json:"total"`
	Items []appsv1.Deployment `json:"items"`
}

// deployment 转 DataCells
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {

	cells := make([]DataCell, len(deployments))

	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

// DataCells 转 deployment
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {

	Deployment := make([]appsv1.Deployment, len(cells))

	for i := range cells {

		Deployment[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}

	return Deployment

}

// 修改deployment的副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {

	//	 获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName,
		metav1.GetOptions{})

	if err != nil {
		logger.Error(errors.New("获取deployment副本数量信息失败," + err.Error()))
		return 0, errors.New("获取deployment副本数量信息失败," + err.Error())
	}
	// 修改副本数量
	scale.Spec.Replicas = int32(scaleNum)
	// 更新副本数，传入scale对象
	newscale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(),
		deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新deployment副本数量失败," + err.Error()))
		return 0, errors.New("更新deployment副本数量失败," + err.Error())
	}
	return newscale.Spec.Replicas, nil

}

// 创建deployment
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {

	//	 组装一个deployment yaml对应的数据类型
	deployment := &appsv1.Deployment{
		//	objectMeta 中定义资源名，命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		//	spec 中定义副本数，选择器，以及pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				// 定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				//	定义容器名，镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},

		//	Status 定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus[]对象
		Status: appsv1.DeploymentStatus{},
	}
	// 判断是否打开健康检查功能，若打开，
	if data.HealthCheck {
		//	设置一个容器的readinessProbe,因为我们pod中只有一个容器，所以直接使用index 0即可
		//	  若pod中有多个容器，则这里需要使用for循环去定义了
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//	 intstr.IntOrstring的作用是端口可以定义为整形，也可以定义为字符串
					//	 Type=0则表示表示该结构体实例内的数据为整形，转json时只用IntVal的数据
					//	 type=1则表示该结构体实例中的数据为字符串，转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//	初始化等待的时间
			InitialDelaySeconds: 5,
			//	定义超时的时间
			TimeoutSeconds: 5,
			//	执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}

		//	定义容器的limit和request资源.  这里的request就是最小需求
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	}
	//	 调用sdk创建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败," + err.Error()))
		return errors.New("创建Deployment失败," + err.Error())
	}
	return nil
}

// 删除deployment
func (d *deployment) DelDeployment(deploymentName, namespace string) (err error) {

	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除deployment失败," + err.Error()))
		errors.New("删除deployment失败," + err.Error())
	}
	return nil

}

// 重启deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {

	// 重启原理： 修改一个deployment yaml文件中的随便一个值，会自动触发deployment的机制重启
	// 我们这里改 annotation

	// patchData map组装数据
	patchData := map[string]any{
		"spec": map[string]any{
			"template": map[string]any{
				"metadata": map[string]any{
					"annotations": []map[string]any{
						{
							"name": "Restart_", "value": strconv.FormatInt(time.Now().Unix(), 10),
						},
					},
				},
			},
		},
	}
	// 序列化为字节，因为patch方法只接收字节类型的参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("json序列化失败," + err.Error()))
		return errors.New("json序列化失败," + err.Error())
	}
	// 调用patch方法更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启deployment失败," + err.Error()))
		errors.New("重启deployment失败," + err.Error())
	}
	return nil
}

// 更新deployment

func (d *deployment) UpdateDeployment(deploymentName, namespace, Content string) error {

	var deployment = &appsv1.Deployment{}

	//反序列化
	err := json.Unmarshal([]byte(Content), deployment)
	if err != nil {
		logger.Error(errors.New("反序列化失败," + err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}

	//	更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {

		logger.Error(errors.New("更新deployment失败," + err.Error()))
		return errors.New("更新deployment失败," + err.Error())

	}
	return nil

}

// 获取deployment信息

func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (data *DeploymentsResp, err error) {
	// context.TODO() 用于声明一个空的上下文，用于List方法内设置这个请求的超时;具体看源码，源码这里配置是为了定义一个超时时间
	// metav1.ListOption 用于过滤List数据，如使用 label , field 等
	deploymentsList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 打日志给自己看，排错使用
		logger.Info("获取Deployment列表失败," + err.Error())
		// 返回给上一层，最终返回给前端，前端打印出这个error
		return nil, errors.New("获取Deployment列表失败," + err.Error())
	}

	//实例化dataSelector结构体,组装数据

	selectableData := &dataSelector{

		GenericDatalist: d.toCells(deploymentsList.Items),
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

	deployments := d.fromCells(dataList.GenericDatalist)

	return &DeploymentsResp{
		Total: total,
		Items: deployments,
	}, nil

}

// 获取deployment详情
func (d *deployment) GetDeploymentsDetail(deploymentName, namespace string) (pod *appsv1.Deployment, err error) {
	deployment, err := K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment详情失败," + err.Error())
		return nil, errors.New("获取deployment失败," + err.Error())
	}

	return deployment, err
}

// 获取每个namespace下的deployment数量

type DeploymentsNp struct {
	Namespace     string
	DeploymentNum int
}

func (d *deployment) GetDeploymentNumberNp() (deploymentsNps []*DeploymentsNp, err error) {

	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		//	获取pod列表
		DeploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//	 组装数据
		deploymentsNp := &DeploymentsNp{
			Namespace:     namespace.Name,
			DeploymentNum: len(DeploymentList.Items),
		}
		//	 添加到podsNps数组中
		deploymentsNps = append(deploymentsNps, deploymentsNp)
	}

	return deploymentsNps, nil
}
