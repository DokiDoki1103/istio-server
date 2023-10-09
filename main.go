package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"istio-server/kubernetes"
	"istio-server/model"
	"istio-server/prometheus"
	"istio-server/router"
	"log"
	"net/http"
	"runtime/debug"
)

func ErrorHandler(c *gin.Context) {
	// 通过 recover 捕获任何发生的 panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack() // 打印堆栈信息
			c.String(http.StatusInternalServerError, "服务器开小差了,请稍后再试")
		}
	}()
	c.Next()
}

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
	err = model.InitDB(k8sClient)
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
	server.Use(ErrorHandler)
	router.SetRouter(server, buildFS, indexPage)
	err = server.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
