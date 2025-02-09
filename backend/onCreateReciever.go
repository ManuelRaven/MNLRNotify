package main

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/exp/rand"
)

func (app *application) useonCreateReciever() {
	app.pb.OnRecordCreate("reciever").BindFunc(func(e *core.RecordEvent) error {
		reciever := e.Record
		if reciever == nil {
			return e.Next()
		}
		if reciever.Get("type") == "gotify" {
			// Generate a gotify compatible token by starting with an A and 14 random characters
			reciever.Set("token", generateRandomString(14))
		}
		return e.Next()
	})
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return "A" + string(b)
}
