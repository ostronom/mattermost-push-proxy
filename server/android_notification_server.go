// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package server

import (
	"fmt"

	"github.com/alexjlockwood/gcm"
	"github.com/kyokomi/emoji"
)

type AndroidNotificationServer struct {
	AndroidPushSettings AndroidPushSettings
}

func NewAndroideNotificationServer(settings AndroidPushSettings) NotificationServer {
	return &AndroidNotificationServer{AndroidPushSettings: settings}
}

func (me *AndroidNotificationServer) Initialize() bool {
	LogInfo(fmt.Sprintf("Initializing Android notificaiton server for type=%v", me.AndroidPushSettings.Type))

	if len(me.AndroidPushSettings.AndroidApiKey) == 0 {
		LogError("Android push notifications not configured.  Mssing AndroidApiKey.")
		return false
	}

	return true
}

func (me *AndroidNotificationServer) SendNotification(msg *PushNotification) PushResponse {
	var data map[string]interface{}
	if len(msg.Message) > 0 {
		data["message"] = emoji.Sprint(msg.Message)
	}
	for k, v := range msg.CustomData {
		data[k] = v
	}

	regIDs := []string{msg.DeviceId}
	gcmMsg := gcm.NewMessage(data, regIDs...)

	sender := &gcm.Sender{ApiKey: me.AndroidPushSettings.AndroidApiKey}

	if len(me.AndroidPushSettings.AndroidApiKey) > 0 {
		LogInfo(fmt.Sprintf("Sending android push notification for type=%v", me.AndroidPushSettings.Type))
		resp, err := sender.Send(gcmMsg, 2)

		if err != nil {
			LogError(fmt.Sprintf("Failed to send GCM push did=%v err=%v type=%v", msg.DeviceId, err, me.AndroidPushSettings.Type))
			return NewErrorPushResponse("unknown transport error")
		}

		if resp.Failure > 0 {
			if len(resp.Results) > 0 && (resp.Results[0].Error == "InvalidRegistration" || resp.Results[0].Error == "NotRegistered") {
				LogInfo(fmt.Sprintf("Android response failure sending remove code: %v type=%v", resp, me.AndroidPushSettings.Type))
				return NewRemovePushResponse()
			} else {
				LogError(fmt.Sprintf("Android response failure: %v type=%v", resp, me.AndroidPushSettings.Type))
				return NewErrorPushResponse("unknown send response error")
			}
		}
	}

	return NewOkPushResponse()
}
