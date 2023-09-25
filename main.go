package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"istio-server/kubernetes"
	"istio-server/model"
	"istio-server/prometheus"
	"istio-server/router"
	"log"
)

//go:embed web/dist
var buildFS embed.FS

//go:embed web/dist/index.html
var indexPage []byte

func main() {
	// Initialize Kubernetes Client
	k8sClient, err := kubernetes.NewClientFromConfig()
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	// Initialize Prometheus Client
	_, err = prometheus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create prometheus client: %s", err)
	}

	// Initialize SQL Database
	pwd, err := k8sClient.GetMySQLPassword()
	if err != nil {
		log.Fatalf("Failed to get MySQL password: %s", err)
	}
	err = model.InitDB(pwd)
	if err != nil {
		log.Fatalf("Failed to create SQL client: %s", err)
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			log.Fatalf("Failed to close SQL client: %s", err)
		}
	}()

	// Initialize Gin Router
	server := gin.Default()
	router.SetRouter(server, buildFS, indexPage)
	err = server.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
