package tablenamespace

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {

	uls := TableListSource{}
	group.GET("", uls.Get)
}
