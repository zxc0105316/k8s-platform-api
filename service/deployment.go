package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Deployment deployment

type deployment struct {

}
type DeployCreate struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas int32 `json:"replicas"`
	Image string `json:"image"`
	Label map[string]string `json:"label"`
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
	ContainerPort int32 `json:"container_port"`
	HealthCheck bool `json:"health_check"`
	HealthPath string `json:"health_pathe"`


}

// deployment 转 DataCells
func (d *deployment) toCells (deployments []appsv1.Deployment)[]DataCell{

	cells := make([]DataCell,len(deployments))

	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}



// DataCells 转 deployment
func (d *deployment)fromCells(cells []DataCell)[]appsv1.Deployment{

	Deployment := make([]appsv1.Deployment,len(cells))

	for i := range cells{

		Deployment[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}

	return Deployment

}


// 修改deployment的副本数
func (d *deployment) ScaleDeployment(deploymentName,namespace string,scaleNum int)(replica int32 ,err error){

//	 获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale,err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(),deploymentName,
		metav1.GetOptions{})

	if err !=nil {
		logger.Error(errors.New("获取deployment副本数量信息失败," + err.Error()))
		return 0,errors.New("获取deployment副本数量信息失败," + err.Error())
	}
	// 修改副本数量
	scale.Spec.Replicas = int32(scaleNum)
	// 更新副本数，传入scale对象
	newscale,err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(),
		deploymentName,scale,metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新deployment副本数量失败," + err.Error()))
		return 0,errors.New("更新deployment副本数量失败," + err.Error())
	}
	return newscale.Spec.Replicas,nil

}

// 创建deployment
func (d *deployment) CreateDeployment(data *DeployCreate) (err error){

//	 组装一个deployment yaml对应的数据类型
	deployment := &appsv1.Deployment{
	//	objectMeta 中定义资源名，命名空间以及标签
		ObjectMeta:metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
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
					Name: data.Name,
					Labels: data.Label,
				},
			//	定义容器名，镜像和端口
				Spec: corev1.PodSpec{
					Containers:

				}
			}
		}
	}
}