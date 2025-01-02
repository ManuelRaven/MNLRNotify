package main

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

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

func (app *application) useOnMessage() {
	// fires only for "users" and "articles" records
	app.pb.OnRecordAfterCreateSuccess("message").BindFunc(func(e *core.RecordEvent) error {
		// e.App
		// e.Record

		var deliveryMessageLog string = ""
		var deliveryState string = "success"

		// Expand the record to include the "channel" field
		channelId := e.Record.Get("channel").(string)
		if channelId == "" {
			log.Println("channelId is empty")
		}

		records, err := app.pb.FindRecordsByFilter(
			"sender",                           // collection
			"channel ~ {:channelid}",           // filter
			"-created",                         // sort
			999999,                             // limit
			0,                                  // offset
			dbx.Params{"channelid": channelId}, // optional filter params
		)
		if err != nil {
			log.Println("error finding records", err)
		}

		for _, record := range records {
			var senderName = record.GetString("name")
			log.Println("Send Message to Sender: " + senderName)
			var sendurl = record.GetString("sendurl")

			// Get splitLimit from sender record
			splitLimit := record.GetInt("splitLimit")
			messageText := e.Record.GetString("text")

			// Split message if necessary
			messageParts := splitMessage(messageText, splitLimit)

			for i, part := range messageParts {

				err := shoutrrr.Send(sendurl, part)
				if err != nil {
					log.Println("error sending message part", i+1, err)
					deliveryMessageLog += senderName + ": Part " + fmt.Sprint(i+1) + " - " + err.Error() + "\n\n"
					deliveryState = "failure"
				} else {
					log.Println("Message part", i+1, "sent successfully")
					deliveryMessageLog += senderName + ": Part " + fmt.Sprint(i+1) + " - Message sent successfully\n\n"
				}
			}

			// Update the deliveryState and deliveryMessageLog
			e.Record.Set("deliveryState", deliveryState)
			e.Record.Set("deliveryMessage", deliveryMessageLog)
			err = app.pb.Save(e.Record)
			if err != nil {
				log.Println("error saving record", err)
			}
		}

		return e.Next()
	})

}
