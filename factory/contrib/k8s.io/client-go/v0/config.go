package client

import "github.com/xgodev/boost/config"

const (
	root              = "boost.factory.k8s"
	tp                = ".type"
	kubeConfigPath    = ".kubeConfigPath"
	kubeConfigContext = ".kubeConfigContext"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+tp, "KUBECONFIG", "defines client type KUBECONFIG/INCLUSTER")
	config.Add(path+kubeConfigPath, "~/.kube/config", "defines kubeconfig request")
	config.Add(path+kubeConfigContext, "", "defines kubeconfig context")
}
