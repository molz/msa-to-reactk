package main

import (
	"encoding/json"
	"log"
	"fmt"
	"time"
)

type PushEvent struct {
	UID         string      `json:"uid"`
	Day         string      `json:"day"`
	CreatedDate string      `json:"created_date"`
	Payload     interface{} `json:"payload"`
	ClientId    string      `json:"client_id"`
	Type        string      `json:"type"`
}

type Event struct {
	ConfigId  string                 `json:"config_id"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"data"`
	SentDate  string                 `json:"sent_date"`
}

func (e *Event) getPushBody() []byte {
	pushEvent := PushEvent{
		UID:         e.getUID(),
		Day:         e.getDay(),
		CreatedDate: e.getCreatedDate(),
		ClientId:    e.getClientId(),
		Type:        e.getType(),
		Payload:     e.Payload,
	}
	b, err := json.Marshal(pushEvent)
	if err != nil {
		failByLabel.WithLabelValues("fail_marshal_push_body").Inc()
		log.Printf("fail to marshal pushEvent: %s", err.Error())
		return nil
	}
	return b
}

func (e *Event) getDay() string {
	if len(e.SentDate) >= 10 {
		return e.SentDate[:10]
	}
	return ""
}
 const (
 	timeFormatISO8601 = "2006-01-02 15:04:05-0700"
 )
// output format : 2018-09-25 06:08:32.922+0000
func (e *Event) getCreatedDate() string {
	t, err := time.Parse("2006-01-02[T]15:04:05[Z]", e.SentDate)
	if err != nil {
		return time.Now().Format(timeFormatISO8601)
	}
	return t.Format(timeFormatISO8601)
}

func (e *Event) getType() string {
	return e.EventType
}

func (e *Event) getPayloadType() string {
	if d, ok := e.Payload["type"]; ok {
		return d.(string)
	}
	return ""
}

func (e *Event) getUID() string {
	if e.getPayloadType() == "User" {
		if d, ok := e.Payload["id_str"]; ok {
			return d.(string)
		}
		return ""
	}
	if d, ok := e.Payload["owner"]; ok {
		b, _ := json.Marshal(d)
		var data struct{ IdStr string `json:"id_str"` }
		json.Unmarshal(b, &data)
		return data.IdStr
	}
	return ""
}

func (e *Event) getClientId() string {
	return fmt.Sprintf("msa_%s", e.ConfigId)
}

func (e *Event) isValid() bool {
	return e.ConfigId != "" && e.getUID() != "" && e.getType() != "" && e.getDay() != ""
}
