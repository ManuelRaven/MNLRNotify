package main

import (
	"log"
	"time"

	"github.com/pocketbase/dbx"
)

func (app *application) useCronCleanupMessages() {
	app.pb.Cron().MustAdd("MessageCleanup", "*/1 * * * *", func() {
		// Get all Channel records. a Channel Record has the field lifetime_seconds. Then get for each channel all messages older than lifetime_seconds and delete them. If the lifetime_seconds is 0, the messages are never deleted.
		// Get all Channel records
		channels, err := app.pb.FindRecordsByFilter("channel", "", "", 999999, 0, nil)
		if err != nil {
			log.Println("error finding channels", err)
		}

		for _, channel := range channels {
			// Get lifetime_seconds
			lifetimeSeconds := channel.GetFloat("lifetime_seconds")

			// Get all messages older than lifetime_seconds
			if lifetimeSeconds > 0 {
				// Use RFC3339Nano format with UTC timezone
				cutoffTime := time.Now().UTC().Add(-time.Duration(lifetimeSeconds) * time.Second).Format("2006-01-02 15:04:05.999Z")

				messages, err := app.pb.FindRecordsByFilter(
					"message",
					"channel = {:channelid} && created < {:created}",
					"",
					999999,
					0,
					dbx.Params{
						"channelid": channel.Id,
						"created":   cutoffTime,
					},
				)
				if err != nil {
					log.Println("error finding messages", err)
				}

				// Delete messages
				for _, message := range messages {
					err := app.pb.Delete(message)
					if err != nil {
						log.Println("error deleting message", err)
					}

					log.Println("CRON: deleted message", message.Id)
				}
			}
		}
	})

}
