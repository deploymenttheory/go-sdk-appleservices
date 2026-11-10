package main

import (
	"context"
	"fmt"
	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/webhooks"
	"log"
	"os"
	"time"
)

func main() {
	c, err := axm.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, _, err := c.AXMAPI.Webhooks.CreateV1(ctx, &webhooks.WebhookCreateRequest{Data: webhooks.WebhookCreateRequestData{
		Type:       webhooks.ResourceTypeWebhooks,
		Attributes: webhooks.WebhookCreateRequestAttributes{Name: "Device notifications", Enabled: true, URL: os.Getenv("WEBHOOK_URL"), Secret: os.Getenv("WEBHOOK_SECRET"), EventTypes: []string{webhooks.EventTypeOrgDeviceAttributeChanged}},
	}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created webhook: %s\n", result.Data.ID)
}
