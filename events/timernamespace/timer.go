package timernamespace

import (
	"io"
	"net/http"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	events "github.com/kataras/go-events"
)

type TimerSource struct {
}

// @Summary 获取用户列表信息
// @Tags user
// @Produce  text/event-stream
// @Param   channelid  path   string   true  "频道id"
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 404 {string} ResultResponse "未找到指定资源"
// @Router /timer/{channelid} [get]
func (s *TimerSource) Get(c *gin.Context) {
	cq := &ListenQuery{}
	err := c.BindUri(cq)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	ssech := make(chan *sse.Event, 10)
	e := events.EventName(cq.ChannelID)
	l := events.Listener(func(payload ...interface{}) {
		for _, pl := range payload {
			switch pl.(type) {
			case *sse.Event:
				{
					ssech <- pl.(*sse.Event)
				}
			default:
				{
					close(ssech)
				}
			}

		}
	})
	PubSub.AddListener(e, l)
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()
	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			{
				ok := PubSub.RemoveListener(e, l)
				if ok {
					close(ssech)
				}
				log.Debug("watching close", log.Dict{
					"channelid": cq.ChannelID,
				})
				return false
			}
		case message, isopen := <-ssech:
			{
				if isopen {
					c.Render(-1, *message)
					return true
				} else {
					return false
				}

			}
		}
	})
	return
}
