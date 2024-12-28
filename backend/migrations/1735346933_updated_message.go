package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_669929365")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != \"\" && @request.body.channel.owner = @request.auth.id && (@request.body.deliveryState:isset = false || @request.body.deliveryState = 'pending') && @request.body.deliveryMessage:isset = false"
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2653034478",
			"max": 0,
			"min": 0,
			"name": "deliveryMessage",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_669929365")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != \"\" && @request.body.channel.owner = @request.auth.id && (@request.body.deliveryState:isset = false || @request.body.deliveryState = 'pending')"
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text2653034478")

		return app.Save(collection)
	})
}
