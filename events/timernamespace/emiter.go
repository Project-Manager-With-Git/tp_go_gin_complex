package timernamespace

import (
	"github.com/Golang-Tools/pubsubmanager"
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

var PubSub = pubsubmanager.NewPubSubManager()
