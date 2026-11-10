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
	id := os.Getenv("WEBHOOK_ID")
	enabled := false // A pointer transmits false; nil leaves the attribute unchanged.
	result, _, err := c.AXMAPI.Webhooks.UpdateByWebhookIDV1(ctx, id, &webhooks.WebhookUpdateRequest{Data: webhooks.WebhookUpdateRequestData{
		Type: webhooks.ResourceTypeWebhooks, ID: id, Attributes: &webhooks.WebhookUpdateRequestAttributes{Enabled: &enabled},
	}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated webhook: %s\n", result.Data.ID)
}
