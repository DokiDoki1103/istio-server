package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	_ "istio-server/docs"
)

func SetRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	SetApiRouter(router)
	setWebRouter(router, buildFS, indexPage)
}
