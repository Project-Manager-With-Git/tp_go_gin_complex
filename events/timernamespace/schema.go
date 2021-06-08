package timernamespace

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

type Event struct {
	Id    string `json:"id,omitempty"`
	Event string `json:"event,omitempty"`
	Data  string `json:"data,omitempty"`
	Retry uint   `json:"retry,omitempty"`
}
