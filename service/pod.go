package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	"k8s-platform-api/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
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

// 获取pod详情
func (p *pod) GetPodsDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	thispod, err := K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pod详情失败," + err.Error())
		return nil, errors.New("获取pod失败," + err.Error())
	}

	return thispod, err
}

// 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {

	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})

	if err != nil {
		logger.Error("删除pod失败," + err.Error())
		return errors.New("删除pod失败," + err.Error())
	}

	return nil

}

// 更新pod
// context 参数是请求中传入的pod对象的json数组,是对象的序列化  marshal
// 注意这个pod对象是corev1.pod类型的
func (p *pod) UpdatePod(podName, namespace, Content string) error {

	var pod = &corev1.Pod{}

	//反序列化
	err := json.Unmarshal([]byte(Content), pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败," + err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}

	//	更新pod
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {

		logger.Error(errors.New("更新pod失败," + err.Error()))
		return errors.New("更新pod失败," + err.Error())

	}
	return nil

}

// 获取pod中的容器名字
func (p *pod) GetPodContainer(podName, namespace string) (container []string, err error) {

	//	 获取pod详情，复用上面的方法
	pod, err := p.GetPodsDetail(podName, namespace)

	if err != nil {
		return nil, err
	}

	//	 从pod对象中获取内部的容器名字
	var containers []string
	for _, container := range pod.Spec.Containers {

		containers = append(containers, container.Name)
	}
	return containers, nil
}

// 获取容器日志

func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//	设置日志的配置，容器名，tail的行数
	linelimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &linelimit,
	}

	//	获取request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	//	 发起request请求，返回一个io.ReadCloser类型（等同于response.body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取PodLog失败," + err.Error()))
		return "", errors.New("获取PodLog失败," + err.Error())
	}
	defer podLogs.Close()
	//	将response body写入到缓冲区,目的是为了转换成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {

		logger.Error(errors.New("复制PodLog失败," + err.Error()))
		return "", errors.New("复制PodLog失败," + err.Error())
	}

	return buf.String(), nil

}

// 获取每个namespace下的pod数量

type PodsNp struct {
	Namespace string
	PodNum    int
}

func (p *pod) GetPodNumberNp() (podsNps []*PodsNp, err error) {

	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		//	获取pod列表
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//	 组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		//	 添加到podsNps数组中
		podsNps = append(podsNps, podsNp)
	}

	return podsNps, nil
}
