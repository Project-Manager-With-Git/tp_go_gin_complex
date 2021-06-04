package timernamespace

import (
	"io"
	"net/http"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

type TimerSource struct {
}

// @Summary 监听指定计时器
// @Tags event timer
// @Produce  text/event-stream
// @Param   channelid  path   string   true  "频道id"
// @Success 200 {array} string "用户信息"  Format(chunked)
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 404 {string} ResultResponse "未找到指定资源"
// @Router /event/timer/{channelid} [get]
func (s *TimerSource) Get(c *gin.Context) {
	cq := &ListenQuery{}
	err := c.BindUri(cq)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}

	servclose, ok := PubSub.CloseNotify(cq.ChannelID)
	if !ok {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: "channelid not found"})
		return
	}
	ssech, closech, closefunc, err := PubSub.RegistListener(cq.ChannelID, 10)
	if err != nil {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: err.Error()})
		return
	}

	c.Header("Connection", "keep-alive")
	c.Writer.Flush()
	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-servclose:
			{
				closefunc()
				log.Debug("service close", log.Dict{
					"channelid": cq.ChannelID,
				})
				return false
			}
		case <-closech:
			{
				log.Debug("losed by func", log.Dict{
					"channelid": cq.ChannelID,
				})
				return false
			}
		case <-clientGone:
			{
				closefunc()
				log.Debug("client close", log.Dict{
					"channelid": cq.ChannelID,
				})
				return false
			}
		case message, isopen := <-ssech:
			{
				if isopen {
					log.Debug("channel open", log.Dict{
						"channelid": cq.ChannelID,
					})
					c.Render(-1, *message.(*sse.Event))
					return true
				} else {
					log.Debug("channel close", log.Dict{
						"channelid": cq.ChannelID,
					})
					return false
				}
			}
		}
	})
}
