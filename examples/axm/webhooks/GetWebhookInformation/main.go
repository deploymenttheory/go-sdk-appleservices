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
	result, _, err := c.AXMAPI.Webhooks.GetByWebhookIDV1(ctx, os.Getenv("WEBHOOK_ID"), &webhooks.RequestQueryOptions{Fields: []string{webhooks.FieldName, webhooks.FieldEnabled, webhooks.FieldURL}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Webhook ID: %s\n", result.Data.ID)
}
