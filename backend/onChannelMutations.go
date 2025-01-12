package main

import (
	"github.com/pocketbase/pocketbase/core"
	"mnlr.de/MNLRNotify/executors"
)

func (app *application) useChannelMutations() {
	app.pb.OnRecordCreate("message").BindFunc(func(e *core.RecordEvent) error {
		channelRecord, err := app.pb.FindFirstRecordByData("channel", "id", e.Record.GetString("channel"))
		if err != nil {
			e.Record.Set("deliveryMessage", err.Error())
			return e.Next()
		}

		mutations := channelRecord.GetStringSlice("mutations")
		if len(mutations) == 0 {
			return e.Next()
		}

		messageText := e.Record.GetString("text")

		// Process mutations in order
		for i := 0; i < len(mutations); i++ {
			mutationRecord, err := app.pb.FindRecordById("channel_mutations", mutations[i])
			if err != nil {
				e.Record.Set("deliveryMessage", err.Error())
				return e.Next()
			}

			if executor := executors.GetExecutor(mutationRecord.GetString("executor")); executor != nil {
				messageText = executor.Execute(messageText, mutationRecord.GetString("script"))
			}
		}

		e.Record.Set("text", messageText)
		return e.Next()
	})
}
