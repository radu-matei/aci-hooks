# aci-hooks - https://radu-matei.com/blog/aci-update/

Environment
-----------


This program requires that the following environment variables are set:

Azure credentials to update container groups

```
AZURE_TENANT_ID: contains your Azure Active Directory tenant ID or domain
AZURE_SUBSCRIPTION_ID: contains your Azure Subscription ID
AZURE_CLIENT_ID: contains your Azure Active Directory Application Client ID
AZURE_CLIENT_SECRET: contains your Azure Active Directory Application Secret
```
Information about which container group to update

```
RESOURCE_GROUP_NAME: the name of the resource group where the container group is deployed
CONTAINER_GROUP_NAME: name of the container group
```