package service

import (
	"k8s-platform-api/dao"
	"k8s-platform-api/model"
)

type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"` // Type判断是否带ingress
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

type workflow struct {
}

//type WorkflowResp struct {
//	Total int `json:"total"`
//	Items []*model.Workflow `json:"items"`
//}

var Workflow workflow

func getIngressName(name string) (IngressName string) {

	return name + "-ing"
}

func getServiceName(name string) (ServiceName string) {
	return name + "-svc"
}

// 新增workflow,workflow对象为落库对象
/*
	workflow 分为三个类型 clusterip nodeport ingress
*/
func (w *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//	判断workflow是不是ingress类型,传入空字符串即可
	//  为了判断是否需要新增ingress
	var ingressName string
	if data.Type == "ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}

	//	 组装workflow 中workflow的单条数据
	workflow := &model.Workflow{

		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}

	//   下面这部分如果要认真做，需要考虑事务处理，就是如果资源后续创建失败的话，需要回滚数据
	//	 调用dao层执行数据库的添加操作
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}

	// 创建k8s资源
	err = createWorkflowRes(data)
	if err != nil {
		return err
	}
	return nil
}

// 封装创建workflow对应的k8s资源，这个是实际创建的对象
// 小写开头的函数,作用域只在当前包中,不支持跨包调用
func createWorkflowRes(data *WorkflowCreate) (err error) {

	//   我们这个融合了service和ingress两个的类型，所以service.type的要做判断，非ingress才对
	var serviceType string

	//	 创建deployment对象
	deployment := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	err = Deployment.CreateDeployment(deployment)
	if err != nil {
		return err
	}
	//	 创建service对象

	if data.Type != "Ingress" {

		serviceType = data.Type

	} else {
		serviceType = "cluster"
	}

	service := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}

	err = Service.CreateService(service)
	if err != nil {
		return err
	}

	//	 创建ingress对象
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}

	return nil

}

// 获取workflow列表
func (w *workflow) GetList(name, namespace string, limit, page int) (data *dao.WorkflowResp, err error) {

	data, err = dao.Workflow.GetWorkflows(name, namespace, limit, page)

	if err != nil {
		return nil, err
	}
	return data, nil
}

// 获取单个workflow
func (w *workflow) GetWorkflow(id int) (data *model.Workflow, err error) {

	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil

}

// 删除workflow
func (w *workflow) DelById(id int) (err error) {

	//	 获取workflow数据
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}

	//	删除k8s中的workflow资源

	err = delworkflowRes(workflow)
	if err != nil {
		return err
	}
	//  删除数据库数据

	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}
	return nil

}

// 封装删除workflow
func delworkflowRes(workflow *model.Workflow) (err error) {

	err = Deployment.DelDeployment(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}
	err = Service.DelService(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}
	if workflow.Type == "Ingress" {

		err = Ingress.DeleteIngresses(workflow.Name, workflow.Namespace)
		if err != nil {
			return err
		}
	}
	return
}
