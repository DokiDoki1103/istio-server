package router

import (
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(efs),
	}
}

func Cache() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=604800") // one week
		c.Next()
	}
}

func setWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(Cache())
	router.Use(static.Serve("/", EmbedFolder(buildFS, "web/dist")))
	router.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
