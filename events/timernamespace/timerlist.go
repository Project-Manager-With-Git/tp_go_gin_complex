package timernamespace

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
	"tp_go_gin_complex/events/setting"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/Golang-Tools/redishelper/proxy"
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
	ctx := context.Background()
	pubsub := proxy.Proxy.Subscribe(ctx, "timer::global")
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
					"channelid": "timer::global",
				})

				return false
			}
		case message, isopen := <-ssech:
			{
				if isopen {
					log.Debug("channel open", log.Dict{
						"channelid": "timer::global",
					})
					msg := Event{}
					err := json.UnmarshalFromString(message.Payload, &msg)
					if err != nil {
						log.Error("UnmarshalFromString error", log.Dict{"Payload": message.Payload})
						return true
					}
					if msg.Event == "EOF" {
						log.Debug("publisher close", log.Dict{
							"channelid": "timer::global",
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
						"channelid": "timer::global",
					})
					return false
				}
			}
		}
	})
}

func sendMsg(evt *Event, es string, withgolbal bool) error {
	msg, err := json.MarshalToString(*evt)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), setting.RedisQueryTimeout)
	defer cancel()
	pipe := proxy.Proxy.TxPipeline()
	pipe.Publish(ctx, "timer::"+es, msg)
	if withgolbal {
		pipe.Publish(ctx, "timer::global", msg)
	}
	_, err = pipe.Exec(ctx)
	return err
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
	ctx, cancel := context.WithTimeout(context.Background(), setting.RedisQueryTimeout)
	defer cancel()
	_, err = proxy.Proxy.SAdd(ctx, "timer::channels", es).Result()
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}

	go func(count int) {
		for i := 0; i < count; i++ {
			err := sendMsg(&Event{
				Id:    es + "::" + strconv.Itoa(i),
				Event: "countdown",
				Data:  strconv.Itoa(count - i),
			}, es, true)
			if err != nil {
				log.Error("sendMsg get error", log.Dict{"err": err.Error()})
			} else {
				log.Info("sendMsg ok")
			}
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
		err := sendMsg(&Event{
			Event: "EOF",
		}, es, false)
		if err != nil {
			log.Error("sendMsg get error", log.Dict{"err": err.Error()})
		} else {
			log.Info("send EOF Msg ok")
		}
		ctx, cancel := context.WithTimeout(context.Background(), setting.RedisQueryTimeout)
		defer cancel()
		_, err = proxy.Proxy.SRem(ctx, "timer::channels", es).Result()
		if err != nil {
			c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
			return
		}
	}(cd.Seconds)
	c.PureJSON(200, &CounterDownResponse{ChannelID: es})
}
