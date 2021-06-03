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
	events "github.com/kataras/go-events"
	uuid "github.com/satori/go.uuid"
)

type TimerListSource struct {
}

// @Summary 创建新用户
// @Tags event timer
// @Accept application/json
// @Produce text/event-stream
// @Success 200 {array} string "sse信息"
// @Router /event/timer [get]
func (s *TimerListSource) Get(c *gin.Context) {
	ssech := make(chan *sse.Event, 10)
	closech := make(chan struct{})
	e := events.EventName("default")
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
					"channelid": "default",
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
					c.Render(-1, *message)
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

// @Summary 创建新用户
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
	e := events.EventName(es)
	ge := events.EventName("default")
	go func(count int) {
		for i := 0; i < count; i++ {
			evt := sse.Event{
				Id:    es + "::" + strconv.Itoa(i),
				Event: "countdown",
				Data:  strconv.Itoa(count - i),
			}
			PubSub.Emit(e, &evt)
			PubSub.Emit(ge, &evt)
			time.Sleep(time.Second)
		}
		PubSub.Emit(e, &struct{}{})
		time.Sleep(time.Second)
		PubSub.RemoveAllListeners(e)
	}(cd.Seconds)
	c.PureJSON(200, &CounterDownResponse{ChannelID: es})
}
