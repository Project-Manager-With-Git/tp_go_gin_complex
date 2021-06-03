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

var PubSub = events.New()
