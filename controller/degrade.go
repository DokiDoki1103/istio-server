package controller

import (
	"github.com/gin-gonic/gin"
	"istio-server/kubernetes"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

func GetDegradeRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	config, err := client.GetDegradeConfig(c.Param("ns"), c.Param("name"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "熔断规则查询失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, config.Spec.GetTrafficPolicy().GetOutlierDetection())
}

func PutDegradeRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	var degradeConfig networkingv1alpha3.OutlierDetection
	err := c.ShouldBindJSON(&degradeConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "熔断规则参数错误",
			"error":   err.Error(),
		})
		return
	}

	config, err := client.PutDegradeConfig(c.Param("ns"), c.Param("name"), &degradeConfig)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "熔断规则配置失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, config.Spec.GetTrafficPolicy().GetOutlierDetection())
}

func DelDegradeRule(c *gin.Context) {
	var client = kubernetes.GetK8sClient()

	err := client.DelDegradeConfig(c.Param("ns"), c.Param("name"))
	if err != nil && !errors.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "熔断规则删除失败",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
