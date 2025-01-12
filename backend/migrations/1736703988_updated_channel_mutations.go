package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2100191860")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "select4151479252",
			"maxSelect": 1,
			"name": "executor",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"goja",
				"cmd"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2100191860")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "select4151479252",
			"maxSelect": 1,
			"name": "executor",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"goja"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
