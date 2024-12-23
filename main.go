package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"strings"
	"time"
)

var kubeconfig *string
var namespace *string
var ingressName *string
var iplist *string
var ERR error
var ClientSet *kubernetes.Clientset

func init() {
	// 从命令行获取运行参数
	kubeconfig = flag.String("kubeconfig", "", "Path to the kubeconfig file")
	namespace = flag.String("namespace", "", "Namespace to use")
	ingressName = flag.String("ingressName", "", "Ingress name to update")
	iplist = flag.String("iplist", "", "Comma-separated list of IPs to whitelist")
	flag.Parse()

	// 构建配置
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000)
	if err != nil {
		panic(err.Error())
	}

	// 创建 Kubernetes 客户端
	ClientSet, ERR = kubernetes.NewForConfig(config)

}

func main() {
	nsExist := CheckNamespace(ClientSet, *namespace)
	ips := strings.Split(*iplist, ",")
	if nsExist {
		err := AddIPsToWhitelist(ClientSet, *namespace, *ingressName, ips)
		if err != nil {
			fmt.Println("加白错误：" + err.Error())
			SendToLark("加白错误：" + err.Error())
		}

		// 记录执行结果和时间
		fmt.Println(*namespace, "中的", *ingressName, "IP白名单：", ips, "添加完成", "时间:", time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("错误：传入的Namespace不存在：", *namespace)
		SendToLark("错误：商户" + *namespace + "不存在")
	}

}
