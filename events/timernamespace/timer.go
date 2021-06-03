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
	ssech := make(chan *sse.Event, 10)
	closech := make(chan struct{})
	e := events.EventName(cq.ChannelID)
	// find := false
	EventNames := PubSub.EventNames()
	log.Info("get EventNames", log.Dict{"EventNames": EventNames})
	// if len(EventNames) == 0 {
	// 	find = true
	// } else {
	// 	for _, name := range EventNames {
	// 		if e == name {
	// 			find = true
	// 			break
	// 		}
	// 	}
	// }
	// if !find {
	// 	c.JSON(http.StatusNotFound, &ResultResponse{Message: "channelid not found"})
	// 	return
	// }
	l := events.Listener(func(payload ...interface{}) {
		for _, pl := range payload {
			switch pl := pl.(type) {
			case *sse.Event:
				{
					ssech <- pl
				}
			default:
				{
					close(closech)
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
		case <-closech:
			{
				ok := PubSub.RemoveListener(e, l)
				if ok {
					close(ssech)
				}
				log.Debug("service close", log.Dict{
					"channelid": cq.ChannelID,
				})
				return false
			}
		case <-clientGone:
			{
				ok := PubSub.RemoveListener(e, l)
				if ok {
					close(ssech)
				}
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
					c.Render(-1, *message)
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
