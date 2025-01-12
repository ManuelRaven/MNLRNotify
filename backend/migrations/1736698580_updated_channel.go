package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_866841005")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != \"\" && (@request.body.mutations:isset = false || @request.body.mutations:length = 0 || @request.body.mutations.owner = @request.auth.id)",
			"updateRule": "@request.auth.id != \"\" && @request.auth.id = owner.id && (@request.body.mutations:isset = false || @request.body.mutations:length = 0 || @request.body.mutations.owner = @request.auth.id)"
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_2100191860",
			"hidden": false,
			"id": "relation1335873882",
			"maxSelect": 999,
			"minSelect": 0,
			"name": "mutations",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_866841005")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"createRule": "@request.auth.id != \"\"",
			"updateRule": "@request.auth.id != \"\" && @request.auth.id = owner.id"
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation1335873882")

		return app.Save(collection)
	})
}
