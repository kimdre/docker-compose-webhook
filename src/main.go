package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"net/http"
	"os"
)

const (
	path = "/webhooks"
)

func main() {
	githubWebhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	port := os.Getenv("PORT")

	hook, _ := github.New(github.Options.Secret(githubWebhookSecret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
		if err != nil {
			if errors.Is(err, github.ErrEventNotFound) {
				// ok event wasn't one of the ones asked to be parsed
				fmt.Println("nope")
			}
		}

		fmt.Printf("Received %+v", payload)
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			// DO whatever you want from here
			fmt.Printf("%+v", push)

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	})
	fmt.Println("Server is running on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}
