// fs-dev.go
//go:build !production

package main

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func (app *application) useNtfyEndpoint() {
	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.Any("/ntfy/{topic}", func(re *core.RequestEvent) error {
			// Hello World
			topic := re.Request.PathValue("topic")
			fmt.Println("Hello World", topic)
			re.Response.Write([]byte("Hello World " + topic))
			return nil
		})

		return se.Next()
	})
}
