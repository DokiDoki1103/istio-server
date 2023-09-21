package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"istio-server/kubernetes"
	"istio-server/prometheus"
	"istio-server/router"
	"log"
)

//go:embed web/dist
var buildFS embed.FS

//go:embed web/dist/index.html
var indexPage []byte

func main() {
	_, err := kubernetes.NewClientFromConfig()

	if err != nil {
		fmt.Println("Failed to create istio c" + err.Error())
		log.Fatalf("Failed to create istio client: %s", err)
	}

	_, err = prometheus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create prometheus client: %s", err)
	}

	// 创建http服务
	server := gin.Default()
	router.SetRouter(server, buildFS, indexPage)
	err = server.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
