package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3388177990")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_866841005",
			"hidden": false,
			"id": "relation2734263879",
			"maxSelect": 999,
			"minSelect": 0,
			"name": "channel",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3388177990")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation2734263879")

		return app.Save(collection)
	})
}
