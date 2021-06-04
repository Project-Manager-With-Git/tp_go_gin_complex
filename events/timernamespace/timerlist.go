package timernamespace

import (
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	uuid "github.com/satori/go.uuid"
)

type TimerListSource struct {
}

// @Summary 监听全部计时器
// @Tags event timer
// @Accept application/json
// @Produce text/event-stream
// @Success 200 {array} string "sse信息"
// @Router /event/timer [get]
func (s *TimerListSource) Get(c *gin.Context) {
	ssech, closech, closefunc, err := PubSub.RegistListener("default", 10)
	if err != nil {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: "channelid not found"})
		return
	}
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()
	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-closech:
			{
				log.Debug("closed by func", log.Dict{
					"channelid": "default",
				})
				return false
			}
		case <-clientGone:
			{
				closefunc()
				log.Debug("client close", log.Dict{
					"channelid": "default",
				})
				return false
			}
		case message, isopen := <-ssech:
			{
				if isopen {
					log.Debug("channel open", log.Dict{
						"channelid": "default",
					})
					c.Render(-1, *message.(*sse.Event))
					return true
				} else {
					log.Debug("channel close", log.Dict{
						"channelid": "default",
					})
					return false
				}
			}
		}
	})
}

// @Summary 创建计时器
// @Tags event timer
// @Accept application/json
// @Produce application/json
// @Param counterdown body CounterDownQuery true "倒计时信息"
// @Success 200 {object} CounterDownResponse "倒计时频道信息信息"
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 500 {string} ResultResponse "服务器处理失败"
// @Router /event/timer [post]
func (s *TimerListSource) Post(c *gin.Context) {
	cd := &CounterDownQuery{}
	err := c.ShouldBindBodyWith(cd, binding.JSON)
	if err != nil {
		c.PureJSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	es := uuid.NewV4().String()
	go func(count int) {
		for i := 0; i < count; i++ {
			evt := sse.Event{
				Id:    es + "::" + strconv.Itoa(i),
				Event: "countdown",
				Data:  strconv.Itoa(count - i),
			}
			PubSub.PublishWithDefault(&evt, es)
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
		PubSub.CloseChannel(es)
	}(cd.Seconds)
	c.PureJSON(200, &CounterDownResponse{ChannelID: es})
}
