package events

import (
	"tp_go_gin_complex/events/setting"
	"tp_go_gin_complex/events/timernamespace"

	"github.com/Golang-Tools/redishelper/proxy"
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(app *gin.Engine, redis_url string, redis_query_timeout_ms int) *gin.Engine {

	// 计时器
	timer := app.Group("/v1_0_0/event/timer")
	proxy.Proxy.InitFromURL(redis_url)
	setting.Init(redis_query_timeout_ms)
	timernamespace.Init(timer)
	return app
}
