package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// CheckNamespace 检查namespace是否存在
func CheckNamespace(clientSet *kubernetes.Clientset, namespace string) bool {
	log.Println("正在检查namespace是否存在：", namespace)
	_, err := clientSet.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		log.Println("namespace不存在：", namespace)
		return false
	}
	log.Println("namespace存在：", namespace)
	return true
}

// AddIPsToWhitelist 添加 IP 到白名单
func AddIPsToWhitelist(clientSet *kubernetes.Clientset, namespace, ingressName string, ips []string) error {
	log.Println("正在增加ip到", namespace, "中的", ingressName, "IP为：", ips)
	// 获取 Ingress 对象
	ingress, err := clientSet.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 删除现有的白名单
	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = ""

	// 添加新的 IP 到白名单
	newWhitelist := strings.Join(ips, ",")
	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = newWhitelist

	// 更新 Ingress 对象
	_, err = clientSet.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Println("加白完成:", namespace, "中的", ingressName, "IP为：", ips)

	return nil
}
