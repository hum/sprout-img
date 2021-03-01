package main

import (
	"net/http"

  "github.com/hum/sprout-img"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	/*
	  TODO:
	  Currently only handle for Twitch specific actions

	  Handle event registration, cancellation and restart
	  Create structs to match responses

	  Params:
	    ~/wh?action=[action]&id=[id]
	      • action: type of action to perform on the specific webhook
	        -> add/remove/check
	      • id: user id that identifies the webhook
	*/
	return
}

func main() {
	http.HandleFunc("/", sproutimg.HandleImageCollection)

	/*
	  TODO:
	  Serve as HTTPS -- generate private SSL cert and key -- to be able to handle webhooks

	  http.HandleFunc("/wh", handleWebhook)
	*/

	http.ListenAndServe(":3000", nil)
}
