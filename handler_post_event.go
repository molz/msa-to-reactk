package main

import (
	"github.com/kataras/iris"
	"time"
	"log"
	"net/http"
	"bytes"
)

var (
	httpClient = getHttpClient(10, time.Minute)
	events     = make(chan Event, 10000)
)

func init() {
	for i := 0; i < *concurrentRequest; i++ {
		go func() {
			for event := range events {
				if !event.isValid() {
					log.Printf("invalid payload: %+v", event)
					continue
				}
				pushEvent(&event)
			}
		}()
	}
}

func handlerPostEvent(ctx iris.Context) {
	var event Event
	err := ctx.ReadJSON(&event)
	if err != nil {
		log.Printf("fail to unmarshal json: %s", err.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	events <- event
}

func pushEvent(event *Event) {
	b := event.getPushBody()
	if b == nil {
		log.Printf("fail to get body")
	}
	req, err := http.NewRequest(http.MethodPost, *reactkApi, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("fail to prepare request for event %+v: %s", event, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("fail to push event %s: %s", event, err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		log.Printf("fail to push event, invalid status code %d", resp.StatusCode)
		return
	}
}
