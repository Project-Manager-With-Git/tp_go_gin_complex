package apis

import (
	"github.com/gin-gonic/gin"
)

//RegistPing 注册ping接口到gin
func RegistPing(app *gin.Engine) *gin.Engine {
	app.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return app
}
