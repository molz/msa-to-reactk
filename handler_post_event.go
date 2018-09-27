package main

import (
	"github.com/kataras/iris"
	"time"
	"log"
	"net/http"
	"bytes"
	"strconv"
)

var (
	httpClient = getHttpClient(*concurrentRequest, time.Minute)
	events     = make(chan Event, 10000)
)

func init() {
	for i := 0; i < *concurrentRequest; i++ {
		go func() {
			for event := range events {
				if !event.isValid() {
					failByLabel.WithLabelValues("invalid_payload").Inc()
					log.Printf("invalid payload: %+v", event)
					continue
				}
				pushEvent(&event)
			}
		}()
	}
}

func handlerPostEvent(ctx iris.Context) {
	incomingTotalRequest.Inc()
	var event Event
	err := ctx.ReadJSON(&event)
	if err != nil {
		failByLabel.WithLabelValues("incoming_body_unmarshal").Inc()
		log.Printf("fail to unmarshal json: %s", err.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	requestByConfigId.WithLabelValues(event.ConfigId).Inc()
	events <- event
}

func pushEvent(event *Event) {
	b := event.getPushBody()
	if b == nil {
		log.Printf("fail to get body")
	}
	req, err := http.NewRequest(http.MethodPost, *reactkApi, bytes.NewBuffer(b))
	if err != nil {
		failByLabel.WithLabelValues("prepare_request_fail").Inc()
		reactkResponseCode.WithLabelValues("0").Inc()
		log.Printf("fail to prepare request for event %+v: %s", event, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		failByLabel.WithLabelValues("push_request_error").Inc()
		reactkResponseCode.WithLabelValues("0").Inc()
		log.Printf("fail to push event %s: %s", event, err.Error())
		return
	}
	reactkResponseCode.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		failByLabel.WithLabelValues("invalid_status_code").Inc()
		log.Printf("fail to push event, invalid status code %d", resp.StatusCode)
		return
	}
}
