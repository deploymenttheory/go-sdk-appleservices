package main

import (
	"context"
	"fmt"
	"github.com/deploymenttheory/go-sdk-appleservices/axm"
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
	resp, err := c.AXMAPI.Webhooks.DeleteByWebhookIDV1(ctx, os.Getenv("WEBHOOK_ID"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted webhook (HTTP %d)\n", resp.StatusCode())
}
