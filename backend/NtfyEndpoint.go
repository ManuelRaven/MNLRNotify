package main

import (
	"encoding/json"
	"io"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func (app *application) useNtfyEndpoint() {
	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.POST("/ntfy/{token}", func(re *core.RequestEvent) error {
			// Helper function to send problem response
			sendProblem := func(status int, title, detail string) error {
				re.Response.Header().Set("Content-Type", "application/problem+json")
				re.Response.WriteHeader(status)
				problem := Problem{
					Type:     "about:blank",
					Title:    title,
					Status:   status,
					Detail:   detail,
					Instance: "/ntfy",
				}
				return json.NewEncoder(re.Response).Encode(problem)
			}

			// Get token from path
			token := re.Request.PathValue("token")
			if token == "" {
				return sendProblem(401, "Unauthorized", "No valid authentication token provided")
			}

			// Read message from body
			bodyBytes, err := io.ReadAll(re.Request.Body)
			if err != nil {
				return sendProblem(400, "Bad Request", "Failed to read message body")
			}
			messageText := string(bodyBytes)
			if messageText == "" {
				return sendProblem(400, "Bad Request", "Message body is required")
			}

			// Try to find a receiver with the token
			receiver, err := app.pb.FindRecordsByFilter(
				"reciever",
				"token = {:token} && type = 'ntfy'",
				"-created",
				1,
				0,
				dbx.Params{"token": token},
			)
			if err != nil {
				return sendProblem(500, "Internal Server Error", "Failed to query receiver")
			}

			if receiver == nil {
				return sendProblem(404, "Not Found", "No receiver found for provided token")
			}

			messageCollection, err := app.pb.FindCollectionByNameOrId("message")
			if err != nil {
				return sendProblem(500, "Internal Server Error", "Failed to find message collection")
			}

			// Send message to each channel
			channels := receiver[0].GetStringSlice("channels")
			for _, channel := range channels {
				message := core.NewRecord(messageCollection)
				message.Set("text", messageText)
				message.Set("channel", channel)
				err := app.pb.Save(message)
				if err != nil {
					return sendProblem(500, "Internal Server Error", "Failed to save message")
				}
			}

			// Success response
			re.Response.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(re.Response).Encode(map[string]interface{}{
				"success": true,
				"message": "Messages sent successfully",
			})
		})

		return se.Next()
	})
}
