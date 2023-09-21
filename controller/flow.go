package controller

import (
	"github.com/gin-gonic/gin"
	"istio-server/kubernetes"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

// GetTcpFlowRule 获取tcp限流规则
func GetTcpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	config, err := client.GetTcpFlowConfig(c.Param("ns"), c.Param("name"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "tcp 限流规则查询失败",
			"error":   err.Error(),
		})
		return
	}
	r := config.Spec.GetTrafficPolicy().GetConnectionPool().GetTcp()
	c.JSON(http.StatusOK, r)
}

// PutTcpFlowRule 设置tcp限流规则
func PutTcpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	var tcpFlowConfig networkingv1alpha3.ConnectionPoolSettings_TCPSettings
	err := c.ShouldBindJSON(&tcpFlowConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "tcp 限流规则参数错误",
			"error":   err.Error(),
		})
		return
	}

	config, err := client.PutTcpFlowConfig(c.Param("ns"), c.Param("name"), &tcpFlowConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "tcp 限流规则配置失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, config.Spec.GetTrafficPolicy().GetConnectionPool().GetTcp())
}

// DelTcpFlowRule 删除tcp限流规则
func DelTcpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	err := client.DelTcpFlowConfig(c.Param("ns"), c.Param("name"))
	if err != nil && !errors.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "tcp 限流规则删除失败",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetHttpFlowRule 获取http限流规则
func GetHttpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()
	config, err := client.GetHttpFlowConfig(c.Param("ns"), c.Param("name"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "http 限流规则查询失败",
			"error":   err.Error(),
		})
		return
	}
	r := config.Spec.GetTrafficPolicy().GetConnectionPool().GetHttp()
	c.JSON(http.StatusOK, r)
}

// PutHttpFlowRule 设置http限流规则
func PutHttpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	var httpFlowConfig networkingv1alpha3.ConnectionPoolSettings_HTTPSettings
	err := c.ShouldBindJSON(&httpFlowConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "http 限流参数错误",
			"error":   err.Error(),
		})
		return
	}

	config, err := client.PutHttpFlowConfig(c.Param("ns"), c.Param("name"), &httpFlowConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "http 限流参数配置失败",
			"error":   err.Error(),
		})
		return
	}
	r := config.Spec.GetTrafficPolicy().GetConnectionPool().GetHttp()

	c.JSON(http.StatusOK, r)
}

// DelHttpFlowRule 删除http限流规则
func DelHttpFlowRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()
	err := client.DelHttpFlowConfig(c.Param("ns"), c.Param("name"))
	if err != nil && !errors.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "http 限流规则删除失败",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
