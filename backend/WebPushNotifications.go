package main

import (
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/pocketbase/core"
)

func (app *application) useWebPushNotifications() {

	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// register "GET /hello/{name}" route (allowed for everyone)
		// On Start check VAPID keys exist or create them
		record, _ := app.pb.FindFirstRecordByData("ApplicationKV", "key", "vapid")

		if record == nil {
			// Create VAPID keys
			privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
			if err != nil {
				// TODO: Handle error
				return err
			}

			// Save VAPID keys as json in field value
			vapidKeys := map[string]string{
				"publicKey":  publicKey,
				"privateKey": privateKey,
			}
			// as json string
			vapidKeysJson, err := json.Marshal(vapidKeys)
			if err != nil {
				return err
			}
			collection, err := app.pb.FindCollectionByNameOrId("ApplicationKV")
			if err != nil {
				return err
			}
			record := core.NewRecord(collection)
			record.Set("key", "vapid")
			record.Set("value", string(vapidKeysJson))
			err = app.pb.Save(record)
		}
		return se.Next()
	})

}
