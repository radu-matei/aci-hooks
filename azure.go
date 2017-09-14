package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	tenantID       = getEnvVarOrExit("AZURE_TENANT_ID")
	clientID       = getEnvVarOrExit("AZURE_CLIENT_ID")
	clientSecret   = getEnvVarOrExit("AZURE_CLIENT_SECRET")

	defaultActiveDirectoryEndpoint = azure.PublicCloud.ActiveDirectoryEndpoint
	defaultResourceManagerEndpoint = azure.PublicCloud.ResourceManagerEndpoint
)

func getContainerGroupsClient() (containerinstance.ContainerGroupsClient, error) {
	var containerGroupsClient containerinstance.ContainerGroupsClient

	oAuthConfig, err := adal.NewOAuthConfig(defaultActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return containerGroupsClient, fmt.Errorf("cannot get oAuth configuration: %v", err)
	}

	token, err := adal.NewServicePrincipalToken(*oAuthConfig, clientID, clientSecret, defaultResourceManagerEndpoint)
	if err != nil {
		return containerGroupsClient, fmt.Errorf("cannot get service principal token: %v", err)
	}

	containerGroupsClient = containerinstance.NewContainerGroupsClient(subscriptionID)
	containerGroupsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return containerGroupsClient, nil
}

func updateAzureContainer(resourceGroupName, containerGroupName string, webhookData WebhookData) {
	containerGroupsClient, err := getContainerGroupsClient()
	if err != nil {
		fmt.Printf("cannot get container groups client: %v", err)
	}

	containerGroup, err := containerGroupsClient.Get(resourceGroupName, containerGroupName)
	if err != nil {
		fmt.Printf("cannot get container group: %v", err)
	}

	containers := *containerGroup.Containers
	for index, container := range containers {

		image := *container.Image
		if image == fmt.Sprintf("%s/%s", webhookData.Repository.Namespace, webhookData.Repository.Name) {
			fmt.Println("found image")
			updatedContainer := (*containerGroup.Containers)[index]
			*updatedContainer.Image = fmt.Sprintf("%s/%s:%s", webhookData.Repository.Namespace, webhookData.Repository.Name, webhookData.PushData.Tag)

			_, err := containerGroupsClient.CreateOrUpdate(resourceGroupName, containerGroupName, containerGroup)
			if err != nil {
				fmt.Printf("%v", err)
			}
			fmt.Println("updated image")
		}
		return
	}
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
