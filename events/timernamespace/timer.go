package timernamespace

import (
	"context"
	"io"
	"net/http"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/Golang-Tools/redishelper/proxy"
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
	ctx := context.Background()
	ok, err := proxy.Proxy.SIsMember(ctx, "timer::channels", cq.ChannelID).Result()
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	if !ok {
		c.PureJSON(http.StatusNotFound, &ResultResponse{Message: "channel not in use"})
		return
	}
	pubsub := proxy.Proxy.Subscribe(ctx, "timer::"+cq.ChannelID)
	ssech := pubsub.Channel()
	defer func() {
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
	}()

	c.Header("Connection", "keep-alive")
	c.Writer.Flush()
	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			{
				log.Debug("client close", log.Dict{
					"channelid": "timer::" + cq.ChannelID,
				})
				return false
			}
		case message, isopen := <-ssech:
			{
				if isopen {
					log.Debug("channel open", log.Dict{
						"channelid": "timer::" + cq.ChannelID,
					})
					msg := Event{}
					err := json.UnmarshalFromString(message.Payload, &msg)
					if err != nil {
						log.Error("UnmarshalFromString error", log.Dict{"Payload": message.Payload})
						return true
					}
					if msg.Event == "EOF" {
						log.Debug("publisher close", log.Dict{
							"channelid": "timer::" + cq.ChannelID,
						})
						return false
					}
					c.Render(-1, &sse.Event{
						Id:    msg.Id,
						Event: msg.Event,
						Data:  msg.Data,
						Retry: msg.Retry,
					})
					return true
				} else {
					log.Debug("channel close", log.Dict{
						"channelid": "timer::" + cq.ChannelID,
					})
					return false
				}
			}
		}
	})
}
