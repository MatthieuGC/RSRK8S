package k8s

import (
	"bytes"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/Agorize/EPMToolsBackend/pkg/locate"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// InitClient TO DESCRIBE
func InitClient() (*kubernetes.Clientset, error) {
	home := locate.HomeDir()

	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset, nil
}

// GetNamespaceNames TO DESCRIBE
func GetNamespaceNames(c *kubernetes.Clientset) ([]string, error) {
	namespaces, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var namespaceNames []string
	for _, namespace := range namespaces.Items {
		namespaceNames = append(namespaceNames, namespace.GetName())
	}

	return namespaceNames, err
}

// ExecCmd TO DESCRIBE
func ExecCmd(c *kubernetes.Clientset, pod v1.Pod, cmd []string) (string, string, error) {
	req := c.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	home := locate.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err.Error())
	}

	streamIn := &bytes.Buffer{}
	streamErr := &bytes.Buffer{}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: streamIn,
		Stderr: streamErr,
	})
	if err != nil {
		fmt.Printf("Failed executing command %s on %v/%v", cmd, pod.Namespace, pod.Name)
		panic(err.Error())
	}

	return streamIn.String(), streamErr.String(), nil
}
