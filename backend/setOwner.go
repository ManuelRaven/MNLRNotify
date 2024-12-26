package main

import "github.com/pocketbase/pocketbase/core"

func (app *application) useSetOwner() {
	// fires for every record
	// fires for every collection
	app.pb.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...

		if e.Auth.IsSuperuser() {
			return e.Next()
		}

		// When record has owner field, set it to the current user
		TheField := e.Record.Collection().Fields.GetByName("owner")
		if TheField != nil {
			e.Record.Set("owner", e.Auth.Id)
		}

		return e.Next()
	})

	// fires for every collection
	app.pb.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		// e.App
		// e.Collection
		// e.Record
		// and all RequestEvent fields...

		if e.Auth.IsSuperuser() {
			return e.Next()
		}

		// When record has owner field, set it to the current user
		TheField := e.Record.Collection().Fields.GetByName("owner")
		if TheField != nil {
			e.Record.Set("owner", e.Auth.Id)
		}

		return e.Next()
	})

}
