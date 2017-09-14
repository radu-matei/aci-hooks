package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"
)

var (
	resourceGroupName  = getEnvVarOrExit("RESOURCE_GROUP_NAME")
	containerGroupName = getEnvVarOrExit("CONTAINER_GROUP_NAME")
)

// WebhookData struct is the data sent from Docker Hub webhooks
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
	http.HandleFunc("/", handleWebhook)
	http.ListenAndServe(":8080", nil)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var webhookData WebhookData
	err := decoder.Decode(&webhookData)
	if err != nil {
		log.Printf("cannot decode request body: %v", err)
	}
	defer r.Body.Close()

	fmt.Println("received webhook")
	updateAzureContainer(resourceGroupName, containerGroupName, webhookData)
}
