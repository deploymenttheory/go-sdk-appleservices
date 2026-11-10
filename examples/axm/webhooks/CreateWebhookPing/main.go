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
	result, _, err := c.AXMAPI.Webhooks.CreatePingV1(ctx, &webhooks.WebhookPingCreateRequest{Data: webhooks.WebhookPingCreateRequestData{
		Type:          webhooks.ResourceTypeWebhookPings,
		Relationships: webhooks.WebhookPingCreateRelationships{Webhook: webhooks.ResourceRelationship{Data: webhooks.ResourceIdentifier{Type: webhooks.ResourceTypeWebhooks, ID: os.Getenv("WEBHOOK_ID")}}},
	}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Requested ping: %s\n", result.Data.ID)
}
