// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package server

import (
	"encoding/json"
	"io"
	"log"
)

type PushNotification struct {
	Platform         string                 `json:"platform"`
	DeviceId         string                 `json:"device_id"`
	Sound            string                 `json:"sound"`
	Message          string                 `json:"message"`
	Badge            int                    `json:"badge"`
	ContentAvailable int                    `json:"cont_avavilable"`
	TtlSeconds       int                    `json:"ttl_seconds"`
	CollapseKey      string                 `json:"collapse_key"`
	CustomData       map[string]interface{} `json:"custom"`
}

func (me *PushNotification) ToJson() string {
	b, err := json.Marshal(me)
	if err != nil {
		log.Print("Error marshalling json:", err)
		return ""
	} else {
		return string(b)
	}
}

func PushNotificationFromJson(data io.Reader) *PushNotification {
	decoder := json.NewDecoder(data)
	var me PushNotification
	err := decoder.Decode(&me)
	if err == nil {
		return &me
	} else {
		log.Print("Error unmarshalling json:", err)
		return nil
	}
}
