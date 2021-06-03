package timernamespace

import "github.com/gin-gonic/gin"

func Init(group *gin.RouterGroup) {
	us := TimerSource{}
	uls := TimerListSource{}
	group.GET("", uls.Get)
	group.POST("", uls.Post)
	group.GET("/:channelid", us.Get)
}
