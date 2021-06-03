package events

import (
	"tp_go_gin_complex/events/timernamespace"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(app *gin.Engine) *gin.Engine {

	// 计时器
	timer := app.Group("/v1_0_0/event/timer")
	timernamespace.Init(timer)
	return app
}
