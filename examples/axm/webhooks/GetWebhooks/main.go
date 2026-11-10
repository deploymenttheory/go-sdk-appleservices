package main

import (
	"context"
	"fmt"
	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/webhooks"
	"log"
	"time"
)

func main() {
	c, err := axm.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, _, err := c.AXMAPI.Webhooks.GetV1(ctx, &webhooks.RequestQueryOptions{Fields: []string{webhooks.FieldName, webhooks.FieldEnabled}, Limit: 100})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d webhooks (all pages)\n", len(result.Data))
}
