package main

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Add problem response type for error handling
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (app *application) useGotifyEndpoint() {
	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.POST("/message", func(re *core.RequestEvent) error {
			// Helper function to send problem response
			sendProblem := func(status int, title, detail string) error {
				re.Response.Header().Set("Content-Type", "application/problem+json")
				re.Response.WriteHeader(status)
				problem := Problem{
					Type:     "about:blank",
					Title:    title,
					Status:   status,
					Detail:   detail,
					Instance: "/message",
				}
				return json.NewEncoder(re.Response).Encode(problem)
			}

			// Hello World
			tokenHeader := re.Request.Header.Get("X-Gotify-Key")
			tokenQuery := re.Request.URL.Query().Get("token")
			BearerToken := re.Request.Header.Get("Authorization")

			token := ""
			if tokenHeader != "" {
				token = tokenHeader
			} else if tokenQuery != "" {
				token = tokenQuery
			} else if BearerToken != "" {
				// Remove "Bearer " prefix
				token = BearerToken[7:]
			}

			// When token is empty, return 401
			if token == "" {
				return sendProblem(401, "Unauthorized", "No valid authentication token provided")
			}

			// Try to find a reciever with the token
			reciever, err := app.pb.FindRecordsByFilter(
				"reciever",                 // collection
				"token = {:token}",         // filter
				"-created",                 // sort
				1,                          // limit
				0,                          // offset
				dbx.Params{"token": token}, // optional filter params
			)
			if err != nil {
				return sendProblem(500, "Internal Server Error", "Failed to query receiver")
			}

			// When reciever is not found, return 404
			if reciever == nil {
				return sendProblem(404, "Not Found", "No receiver found for provided token")
			}
			messageCollection, err := app.pb.FindCollectionByNameOrId("message")
			if err != nil {
				return sendProblem(500, "Internal Server Error", "Failed to find message collection")
			}
			// Get channels of reciever and send message to each channel
			// The reciever has the field channels which is a list of channel ids
			channels := reciever[0].GetStringSlice("channels")
			for _, channel := range channels {
				// Send message to channel
				// The message has the fields text and created

				info, err2 := re.RequestInfo()
				messagedata, ok := info.Body["message"].(string)
				if err2 != nil || !ok {
					return sendProblem(400, "Bad Request", "Message is required")
				}

				message := core.NewRecord(messageCollection)
				message.Set("text", messagedata)
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
