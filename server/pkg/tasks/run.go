package tasks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	v1 "k8s.io/api/core/v1"
)

// Run TO DESCRIBE
func Run(task string, pod v1.Pod, params []string) (string, string, error) {
	var formattedParams string
	for _, param := range params {
		formattedParams += fmt.Sprintf(" '%s'", param)
	}

	taskPath := fmt.Sprintf("cmd/production/epmtools/rubytasks/%s.rb", task)
	fileData, err := ioutil.ReadFile(taskPath)
	if err != nil {
		log.Fatal(err)
	}

	// ExecCmdLegacy
	kubeCmd := fmt.Sprintf("kubectl exec -ti %s rails r '%s' %s", pod.Name, fileData, formattedParams)
	cmd := exec.Command("sh", "-c", kubeCmd)

	res, err := cmd.CombinedOutput()
	fmt.Printf(">>>>> %s <<<<<", res)

	if err != nil {
		log.Fatal(err)
	}

	return string(res), "", nil

	// KEEP!

	// cmd := []string{
	// 	"/bin/sh",
	// 	"-c",
	// 	fmt.Sprintf("rails r '%s' %s", string(fileData), formattedParams),
	// }

	// k8sClient, err := k8s.InitClient()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// return k8s.ExecCmd(k8sClient, pod, cmd)
}
