package apis

import (
	"tp_go_gin_complex/apis/usernamespace"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(app *gin.Engine) *gin.Engine {

	// 注册api路由
	// 用户信息
	user := app.Group("/api/v1/user")
	usernamespace.Init(user)
	return app
}
