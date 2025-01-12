package main

import (
	"encoding/json"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func (app *application) useWebPushNotifications() {

	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		shouldReturn, err := createVapidIfMissing(app)
		if shouldReturn {
			return err
		}

		se.Router.GET("/eapi/vapid", func(e *core.RequestEvent) error {
			// Return VAPID public key
			record, _ := app.pb.FindFirstRecordByData("ApplicationKV", "key", "vapid")
			vapidKeys := record.GetString("value")
			// Only return the public key
			var vapidKeysMap map[string]string
			json.Unmarshal([]byte(vapidKeys), &vapidKeysMap)
			return e.JSON(http.StatusOK, vapidKeysMap["publicKey"])
		}).Bind(apis.RequireAuth())

		return se.Next()
	})

}

func createVapidIfMissing(app *application) (bool, error) {
	record, _ := app.pb.FindFirstRecordByData("ApplicationKV", "key", "vapid")

	if record == nil {
		// Create VAPID keys
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			// TODO: Handle error
			return true, err
		}

		// Save VAPID keys as json in field value
		vapidKeys := map[string]string{
			"publicKey":  publicKey,
			"privateKey": privateKey,
		}
		// as json string
		vapidKeysJson, err := json.Marshal(vapidKeys)
		if err != nil {
			return true, err
		}
		collection, err := app.pb.FindCollectionByNameOrId("ApplicationKV")
		if err != nil {
			return true, err
		}
		record := core.NewRecord(collection)
		record.Set("key", "vapid")
		record.Set("value", string(vapidKeysJson))
		err = app.pb.Save(record)
	}
	return false, nil
}
