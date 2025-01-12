package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/containrrr/shoutrrr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// DeliveryResult represents the result of a message delivery attempt
type DeliveryResult struct {
	State   string
	Message string
}

type WebPushMessage struct {
	ChannelName string `json:"channelName"`
	Message     string `json:"message"`
}

// Append combines another delivery result with the current one
func (dr *DeliveryResult) Append(other DeliveryResult) {
	if other.State == "failure" {
		dr.State = "failure"
	}
	if dr.Message != "" && other.Message != "" {
		dr.Message += "\n---\n"
	}
	dr.Message += other.Message
}

// NewDeliveryResult creates a new success result
func NewDeliveryResult() DeliveryResult {
	return DeliveryResult{
		State:   "success",
		Message: "",
	}
}

// Helper function to split message into chunks
func splitMessage(message string, limit int) []string {
	if limit <= 0 || len(message) <= limit {
		return []string{message}
	}

	var chunks []string
	for len(message) > 0 {
		if len(message) <= limit {
			chunks = append(chunks, message)
			break
		}
		chunks = append(chunks, message[:limit])
		message = message[limit:]
	}
	return chunks
}

func (app *application) sendViaWebPush(channelId string, messageText string) DeliveryResult {
	result := NewDeliveryResult()

	webpushDevicesRecords, err := app.pb.FindRecordsByFilter(
		"webpushDevices",
		"channel ~ {:channelid}",
		"-created",
		999999,
		0,
		dbx.Params{"channelid": channelId},
	)

	if err != nil {
		return DeliveryResult{"failure", "Error finding webpush devices: " + err.Error()}
	}

	if len(webpushDevicesRecords) == 0 {
		return DeliveryResult{"success", "No webpush devices found"}
	}

	channelRecord, _ := app.pb.FindFirstRecordByData("channel", "id", channelId)
	channelName := channelRecord.GetString("name")

	pushMessage := WebPushMessage{
		ChannelName: channelName,
		Message:     messageText,
	}

	jsonMessage, _ := json.Marshal(pushMessage)

	record, _ := app.pb.FindFirstRecordByData("ApplicationKV", "key", "vapid")
	vapidKeys := record.GetString("value")

	var vapidKeysMap map[string]string
	json.Unmarshal([]byte(vapidKeys), &vapidKeysMap)

	for _, record := range webpushDevicesRecords {
		deviceSubscription := record.GetString("subscription")
		deviceName := record.GetString("Name")
		owner := record.GetString("owner")

		ownerRecord, _ := app.pb.FindFirstRecordByData("users", "id", owner)
		ownerMail := ownerRecord.GetString("email")

		fmt.Println("Sending to device: ", deviceName, ownerMail)
		s := &webpush.Subscription{}
		json.Unmarshal([]byte(deviceSubscription), s)
		resp, err := webpush.SendNotification(jsonMessage, s, &webpush.Options{
			Subscriber:      ownerMail,
			VAPIDPublicKey:  vapidKeysMap["publicKey"],
			VAPIDPrivateKey: vapidKeysMap["privateKey"],
			TTL:             60 * 60,
		})

		deviceResult := DeliveryResult{}
		if err != nil {
			deviceResult = DeliveryResult{"failure", deviceName + ": " + err.Error()}
		} else {
			deviceResult = DeliveryResult{"success", deviceName + ": Webpush notification sent successfully"}
			resp.Body.Close()
		}
		result.Append(deviceResult)
	}

	return result
}

// sendViaShoutrrrr handles sender retrieval and message sending for Shoutrrr
func (app *application) sendViaShoutrrrr(channelId string, messageText string) DeliveryResult {
	result := NewDeliveryResult()

	senderRecords, err := app.pb.FindRecordsByFilter(
		"sender",
		"channel ~ {:channelid}",
		"-created",
		999999,
		0,
		dbx.Params{"channelid": channelId},
	)
	if err != nil {
		result.State = "failure"
		result.Message = "Error finding senders: " + err.Error()
		return result
	}

	for _, record := range senderRecords {
		senderName := record.GetString("name")
		sendurl := record.GetString("sendurl")
		splitLimit := record.GetInt("splitLimit")

		messageParts := splitMessage(messageText, splitLimit)

		for i, part := range messageParts {
			err := shoutrrr.Send(sendurl, part)
			if err != nil {
				result.Message += senderName + ": Part " + fmt.Sprint(i+1) + " - " + err.Error() + "\n\n"
				result.State = "failure"
			} else {
				result.Message += senderName + ": Part " + fmt.Sprint(i+1) + " - Message sent successfully\n\n"
			}
		}
	}

	return result
}

// sendMessage sends a message using the specified providers
func (app *application) sendMessage(providers []string, channelId string, message string) DeliveryResult {
	result := NewDeliveryResult()

	for _, provider := range providers {
		var providerResult DeliveryResult
		switch provider {
		case "shoutrrr":
			providerResult = app.sendViaShoutrrrr(channelId, message)
		case "webpush":
			providerResult = app.sendViaWebPush(channelId, message)
		default:
			providerResult = DeliveryResult{"failure", fmt.Sprintf("unknown provider: %s", provider)}
		}
		result.Append(providerResult)
	}

	return result
}

func (app *application) useOnMessage() {
	app.pb.OnRecordAfterCreateSuccess("message").BindFunc(func(e *core.RecordEvent) error {
		channelId := e.Record.Get("channel").(string)
		if channelId == "" {
			log.Println("channelId is empty")
			return e.Next()
		}

		messageText := e.Record.GetString("text")
		// Support multiple providers
		providers := []string{"shoutrrr", "webpush"}
		result := app.sendMessage(providers, channelId, messageText)

		e.Record.Set("deliveryState", result.State)
		e.Record.Set("deliveryMessage", result.Message)
		if err := app.pb.Save(e.Record); err != nil {
			log.Println("error saving record", err)
		}

		return e.Next()
	})
}
