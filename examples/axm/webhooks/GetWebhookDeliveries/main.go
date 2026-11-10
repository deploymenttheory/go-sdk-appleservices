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
	result, _, err := c.AXMAPI.Webhooks.GetDeliveriesByWebhookIDV1(ctx, os.Getenv("WEBHOOK_ID"), &webhooks.DeliveryQueryOptions{
		FilterCreatedDateTimeAtOrAfter: []string{time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)},
		FilterCreatedDateTimeAtBefore:  []string{time.Now().UTC().Format(time.RFC3339)},
		Include:                        []string{webhooks.IncludeEvent}, Fields: []string{webhooks.FieldDeliveryState, webhooks.FieldDeliveryEvent},
		EventFields: []string{webhooks.FieldEventPayload, webhooks.FieldEventIsPingEvent}, FilterDeliveryState: []string{webhooks.DeliveryStateSucceeded}, Limit: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d deliveries and %d included events (all pages)\n", len(result.Data), len(result.Included))
}
