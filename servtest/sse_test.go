package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"tp_go_gin_complex/events/timernamespace"

	log "github.com/Golang-Tools/loggerhelper"
	sse_test "github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:5000/v1_0_0/event/timer", bytes.NewBuffer([]byte(`{"second": 10}`)))
	req.Header.Set("Content-Type", "application/json")
	httpclient := &http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)

	msg := timernamespace.CounterDownResponse{}
	json.Unmarshal(body, &msg)
	log.Info("get res", log.Dict{"msg": msg})
	sse_client := sse_test.NewClient("http://localhost:5000/v1_0_0/event/timer/" + msg.ChannelID)
	sse_client.Subscribe("messages", func(msg *sse_test.Event) {
		// Got some data!
		fmt.Println(string(msg.Data))
	})
}
