package watcher

import (
	"flag"
	"fmt"
	"time"

	"github.com/nandiheath/k8s-node-monitor/internal/config"
	"github.com/nandiheath/k8s-node-monitor/internal/dns"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

type Watcher struct {
}

func New() *Watcher {
	m := Watcher{}
	return &m
}

func (m *Watcher) Start() {

	dnsService := dns.New()
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", config.KubeConfigPath, "absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		nodes, err := clientset.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

		var addresses []string
		for _, v := range nodes.Items {
			for _, addr := range v.Status.Addresses {
				if addr.Type == v12.NodeInternalIP {
					fmt.Printf("internal IP address: %s\n", addr.Address)
					addresses = append(addresses, addr.Address)
				}
			}
		}

		//dnsService.UpdateDNS(addresses)
		dnsService.UpdateDNSV2(addresses, []dns.DNSConfig{
			{
				0,5, 30107,
			},
			{
				10,5, 30207,
			},
			{
				20,5, 30307,
			},
		})
		time.Sleep(5 * time.Minute)
	}
}
