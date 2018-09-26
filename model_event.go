package main

import (
	"encoding/json"
	"log"
	"fmt"
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
	ConfigId string                 `json:"config_id"`
	Payload  map[string]interface{} `json:"data"`
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
		log.Printf("fail to marshal pushEvent: %s", err.Error())
		return nil
	}
	return b
}

func (e *Event) getDay() string {
	if d, ok := e.Payload["created_date"]; ok && len(d.(string)) >= 10 {
		return d.(string)[:10]
	}
	return ""
}

func (e *Event) getCreatedDate() string {
	if d, ok := e.Payload["created_date"]; ok {
		return d.(string)
	}
	return ""
}

func (e *Event) getType() string {
	if d, ok := e.Payload["type"]; ok {
		return d.(string)
	}
	return ""
}

func (e *Event) getUID() string {
	if e.getType() == "User" {
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
