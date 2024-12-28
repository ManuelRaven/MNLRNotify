package main

import (
	"log"

	"github.com/containrrr/shoutrrr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

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
			err := shoutrrr.Send(sendurl, e.Record.GetString("text"))

			if err != nil {
				log.Println("error sending message", err)
				// Append to deliveryMessageLog
				deliveryMessageLog += senderName + ": " + err.Error() + "\n\n"
				deliveryState = "failure"
			} else {
				log.Println("Message sent successfully")
				// Append to deliveryMessageLog
				deliveryMessageLog += senderName + ": " + "Message sent successfully" + "\n\n"
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
