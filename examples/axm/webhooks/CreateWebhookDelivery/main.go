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
	result, _, err := c.AXMAPI.Webhooks.CreateDeliveryV1(ctx, &webhooks.WebhookDeliveryCreateRequest{Data: webhooks.WebhookDeliveryCreateRequestData{
		Type:          webhooks.ResourceTypeWebhookDeliveries,
		Relationships: webhooks.WebhookDeliveryCreateRelationships{Template: webhooks.ResourceRelationship{Data: webhooks.ResourceIdentifier{Type: webhooks.ResourceTypeWebhookDeliveries, ID: os.Getenv("WEBHOOK_DELIVERY_ID")}}},
	}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Requested redelivery: %s\n", result.Data.ID)
}
