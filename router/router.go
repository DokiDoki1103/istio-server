package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"istio-server/controller"
	_ "istio-server/docs"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		apiRouter.GET("/status", controller.GetStatus)

		// metricsRouter 图表监控
		metricsRouter := apiRouter.Group("/metrics/graph")
		{
			metricsRouter.GET("/gen", controller.GraphNamespaces)
			metricsRouter.POST("/count/:name", controller.GetHttpCount)
			metricsRouter.POST("/qps/:name", controller.GetHttpQps)
			metricsRouter.POST("/time/:name", controller.GetHttpTime)
			metricsRouter.POST("/bytes/:name", controller.GetHttpBytes)
		}

		configRouter := apiRouter.Group("/config")
		{
			// flowRouter 限流 http tcp
			flowRouter := configRouter.Group("/flow")
			{
				flowRouter.GET("/http/:ns/:name", controller.GetHttpFlowRule)
				flowRouter.PUT("/http/:ns/:name", controller.PutHttpFlowRule)
				flowRouter.DELETE("/http/:ns/:name", controller.DelHttpFlowRule)

				flowRouter.GET("/tcp/:ns/:name", controller.GetTcpFlowRule)
				flowRouter.PUT("/tcp/:ns/:name", controller.PutTcpFlowRule)
				flowRouter.DELETE("/tcp/:ns/:name", controller.DelTcpFlowRule)
			}

			// degradeRouter 熔断
			degradeRouter := configRouter.Group("/degrade")
			{
				degradeRouter.GET("/:ns/:name", controller.GetDegradeRule)
				degradeRouter.PUT("/:ns/:name", controller.PutDegradeRule)
				degradeRouter.DELETE("/:ns/:name", controller.DelDegradeRule)
			}

		}
	}
}
