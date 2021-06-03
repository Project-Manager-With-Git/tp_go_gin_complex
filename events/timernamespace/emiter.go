package timernamespace

import (
	events "github.com/kataras/go-events"
)

type ListenQuery struct {
	ChannelID string `uri:"channelid" binding:"required"`
}
type ResultResponse struct {
	Succeed bool   `json:"succeed"`
	Message string `json:"message,omitempty"`
}
type CounterDownQuery struct {
	Seconds int `json:"second"`
}

type CounterDownResponse struct {
	ChannelID string `json:"channelid"`
}
type SSEvent struct {
	Event string
	Id    string
	Retry uint
}

var PubSub = events.New()
