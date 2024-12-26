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

		// appendable collection
		var sendurls []string

		for _, record := range records {
			log.Println("Send Message to Sender", record.GetString("name"))
			sendurls = append(sendurls, record.GetString("sendurl"))
		}

		for _, sendurl := range sendurls {
			err := shoutrrr.Send(sendurl, e.Record.GetString("text"))
			if err != nil {
				log.Println("error sending message", err)
			}
		}

		return e.Next()
	})

}
