package middlewares

import (
	log "github.com/Golang-Tools/loggerhelper"

	"github.com/gin-gonic/gin"
)

//UselessMiddlewareFactory 没什么用的中间件
func UselessMiddlewareFactory(name string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		log.Info("Hello", log.Dict{"name": name})
		ctx.Writer.Header().Set("Author", name)
		ctx.Next()
		log.Info("bye", log.Dict{"name": name})
	}
}
