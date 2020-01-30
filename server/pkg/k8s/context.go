package k8s

import (
	"fmt"
	"log"
	"os/exec"
)

// EnsureContext TO DESCRIBE
func EnsureContext(context string) {
	var cmd *exec.Cmd
	if context == "" {
		cmd = exec.Command("kubectx", "agorize-development")
	} else {
		cmd = exec.Command("kubectx", context)
	}

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(stdout))
}
