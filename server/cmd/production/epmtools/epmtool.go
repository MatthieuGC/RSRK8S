package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Agorize/EPMToolsBackend/pkg/k8s"
	"github.com/Agorize/EPMToolsBackend/pkg/tasks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RunTask TO DESCRIBE
type RunTask struct {
	Task     string   `json:"task"`
	Platform string   `json:"platform"`
	Params   []string `json:"params"`
}

func getPlatformsHandler(w http.ResponseWriter, r *http.Request) {
	k8sClient, err := k8s.InitClient()
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := k8s.GetNamespaceNames(k8sClient)
	if err != nil {
		panic(err.Error())
	}

	namespacesJSON, err := json.Marshal(namespaces)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(namespacesJSON))
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("cmd/production/epmtools/rubytasks")
	if err != nil {
		panic(err.Error())
	}

	var tasks []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		taskName := strings.Replace(file.Name(), ".rb", "", 1)
		tasks = append(tasks, taskName)
	}

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(tasksJSON)
}

func getTaskConfig(w http.ResponseWriter, r *http.Request) {
	configPath := fmt.Sprintf("cmd/production/epmtools/rubytasks/%s.json", r.URL.Query()["taskName"][0])

	fileData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(fileData)
}

func runTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}

		var data RunTask
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			panic(err.Error())
		}

		k8sClient, err := k8s.InitClient()
		if err != nil {
			panic(err.Error())
		}

		pods, err := k8sClient.CoreV1().Pods(data.Platform).List(metav1.ListOptions{LabelSelector: "app.kubernetes.io/component=web"})
		if err != nil {
			panic(err.Error())
		}

		stdout, stderr, err := tasks.Run(data.Task, pods.Items[0], data.Params)
		if err != nil {
			panic(err.Error())
		}

		res := regexp.MustCompile(`^.*\s*.*method\.`).ReplaceAllString(stdout, "")
		fmt.Printf("%s ----------- %s", res, stderr)

		w.Write([]byte(res))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	k8s.EnsureContext("")
	http.HandleFunc("/platforms", getPlatformsHandler)
	http.HandleFunc("/tasks", getTasksHandler)
	http.HandleFunc("/task_config", getTaskConfig)
	http.HandleFunc("/run_task", runTaskHandler)
	http.ListenAndServe(":6969", nil)
}
