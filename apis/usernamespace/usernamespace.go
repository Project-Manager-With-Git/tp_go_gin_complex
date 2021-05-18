package usernamespace

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {
	us := UserSource{}
	uls := UserListSource{}
	group.GET("/", uls.Get)
	group.GET("", uls.Get)
	group.POST("/", uls.Post)
	group.POST("", uls.Post)
	group.GET("/:uid", us.Get)
	group.PUT("/:uid", us.Put)
	group.DELETE("/:uid", us.Delete)
}
