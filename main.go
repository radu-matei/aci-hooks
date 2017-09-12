package main

import (
	"encoding/json"
	"net/http"

	"log"
)

// WebhookData struct models the data sent from Docker Hub webhooks
type WebhookData struct {
	PushData struct {
		PushedAt int      `json:"pushed_at"`
		Images   []string `json:"images"`
		Tag      string   `json:"tag"`
		Pusher   string   `json:"pusher"`
	} `json:"push_data"`
	CallbackURL string `json:"callback_url"`
	Repository  struct {
		Status          string `json:"status"`
		Description     string `json:"description"`
		IsTrusted       bool   `json:"is_trusted"`
		FullDescription string `json:"full_description"`
		RepoURL         string `json:"repo_url"`
		Owner           string `json:"owner"`
		IsOfficial      bool   `json:"is_official"`
		IsPrivate       bool   `json:"is_private"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		StarCount       int    `json:"star_count"`
		CommentCount    int    `json:"comment_count"`
		DateCreated     int    `json:"date_created"`
		RepoName        string `json:"repo_name"`
	} `json:"repository"`
}

func main() {
	http.HandleFunc("/", hookHandler)
	http.ListenAndServe(":8080", nil)
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var webhookData WebhookData
	err := decoder.Decode(&webhookData)
	if err != nil {
		log.Printf("cannot decode request body: %v", err)
	}
	defer r.Body.Close()

	log.Printf("Pusher: %s", webhookData.PushData.Pusher)
	log.Printf("%s/%s", webhookData.Repository.Owner, webhookData.Repository.Name)
	log.Printf("Tag: %s", webhookData.PushData.Tag)
}
