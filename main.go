package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 通过kubeconfig 将现有conf文件转换成rest.config 类型的对象
	conf, err := clientcmd.BuildConfigFromFlags("", "G:\\k8s\\config")

	if err != nil {
		panic(err)
	}
	// 根据rest.config 类型的对象，new一个clientset出来,NewForConfig 需要传入一个restful conf指针对象,也就是上面的conf
	clientset, err := kubernetes.NewForConfig(conf)

	if err != nil {
		panic(err)
	}

	// 使用clientset获取pod列表 ， listoption是过滤条件
	podList, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("test")
	for _, pod := range podList.Items {
		fmt.Println(pod.Name, pod.Namespace)

	}

}
