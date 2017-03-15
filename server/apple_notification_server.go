// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package server

import (
	"fmt"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"github.com/Masterminds/glide/msg"
)

type AppleNotificationServer struct {
	ApplePushSettings ApplePushSettings
	AppleClient       *apns.Client
}

func NewAppleNotificationServer(settings ApplePushSettings) NotificationServer {
	return &AppleNotificationServer{ApplePushSettings: settings}
}

func (me *AppleNotificationServer) Initialize() bool {
	LogInfo(fmt.Sprintf("Initializing apple notificaiton server for type=%v", me.ApplePushSettings.Type))

	if len(me.ApplePushSettings.ApplePushCertPrivate) > 0 {
		appleCert, appleCertErr := certificate.FromPemFile(me.ApplePushSettings.ApplePushCertPrivate, me.ApplePushSettings.ApplePushCertPassword)
		if appleCertErr != nil {
			LogCritical(fmt.Sprintf("Failed to load the apple pem cert err=%v for type=%v", appleCertErr, me.ApplePushSettings.Type))
			return false
		}

		if me.ApplePushSettings.ApplePushUseDevelopment {
			me.AppleClient = apns.NewClient(appleCert).Development()
		} else {
			me.AppleClient = apns.NewClient(appleCert).Production()
		}

		return true
	} else {
		LogError(fmt.Sprintf("Apple push notifications not configured.  Mssing ApplePushCertPrivate. for type=%v", me.ApplePushSettings.Type))
		return false
	}
}

func (me *AppleNotificationServer) SendNotification(msg *PushNotification) PushResponse {
	notification := &apns.Notification{}
	notification.DeviceToken = msg.DeviceId
	msgPayload := payload.NewPayload()
	notification.Payload = msgPayload
	notification.Topic = me.ApplePushSettings.ApplePushTopic
	msgPayload.Badge(msg.Badge)
	for k, v := range msg.CustomData{
		msgPayload.Custom(k, v)
	}

	if me.AppleClient != nil {
		LogInfo(fmt.Sprintf("Sending apple push notification type=%v", me.ApplePushSettings.Type))
		res, err := me.AppleClient.Push(notification)
		if err != nil {
			LogError(fmt.Sprintf("Failed to send apple push did=%v err=%v type=%v", msg.DeviceId, err, me.ApplePushSettings.Type))
			return NewErrorPushResponse("unknown transport error")
		}

		if !res.Sent() {
			if res.Reason == "BadDeviceToken" || res.Reason == "Unregistered" || res.Reason == "MissingDeviceToken" {
				LogInfo(fmt.Sprintf("Failed to send apple push sending remove code res ApnsID=%v reason=%v code=%v type=%v", res.ApnsID, res.Reason, res.StatusCode, me.ApplePushSettings.Type))
				return NewRemovePushResponse()

			} else {
				LogError(fmt.Sprintf("Failed to send apple push with res ApnsID=%v reason=%v code=%v type=%v", res.ApnsID, res.Reason, res.StatusCode, me.ApplePushSettings.Type))
				return NewErrorPushResponse("unknown send response error")
			}
		}
	}

	return NewOkPushResponse()
}
