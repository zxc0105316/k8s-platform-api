package service

import (
	"context"
	"fmt"
	"github.com/wonderivan/logger"
	//"google.golang.org/appengine/log"
	"k8s-platform-api/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 用于初始化k8s clientset

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {

	// 通过kubeconfig 将现有conf文件转换成rest.config 类型的对象
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)

	if err != nil {
		panic("获取 k8s client 配置失败," + err.Error())
	} else {
		logger.Info("k8s client 初始化成功")
	}
	// 根据rest.config 类型的对象，new一个clientset出来,NewForConfig 需要传入一个restful conf指针对象,也就是上面的conf
	clientset, err := kubernetes.NewForConfig(conf)

	if err != nil {
		panic("kubernetes clientset对象初始化失败 ," + err.Error())
	} else {
		logger.Info("k8s clientset初始化成功")
	}

	// 使用clientset获取pod列表 ， listoption是过滤条件
	podList, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("clientset 获取pod列表失败," + err.Error())
		logger.Info("k8s clientset 获取pods列表成功")
	}

	fmt.Println("test")
	for _, pod := range podList.Items {
		fmt.Println(pod.Name, pod.Namespace, pod.Spec)

	}
	k.ClientSet = clientset
}
