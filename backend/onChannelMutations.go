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
		for i := range mutations {
			mutationRecord, err := app.pb.FindRecordById("channel_mutations", mutations[i])
			if err != nil {
				e.Record.Set("deliveryMessage", err.Error())
				return e.Next()
			}

			if executor := executors.GetExecutor(mutationRecord.GetString("executor")); executor != nil {
				mep := &executors.MessageExecutorParameters{
					Message: messageText,
					Script:  mutationRecord.GetString("script"),
					PB:      app.pb,
					OwnerID: channelRecord.GetString("owner"),
				}
				messageText = executor.Execute(mep)
			}
		}

		e.Record.Set("text", messageText)
		return e.Next()
	})
}
