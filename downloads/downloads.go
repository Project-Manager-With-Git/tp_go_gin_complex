package downloads

import (
	"tp_go_gin_complex/downloads/tablenamespace"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(app *gin.Engine) *gin.Engine {

	// 注册api路由
	// 用户信息
	user := app.Group("/v1_0_0/download/table")
	tablenamespace.Init(user)
	return app
}
