// Set APPLE_KEY_ID, APPLE_ISSUER_ID and APPLE_PRIVATE_KEY_PATH (or
// APPLE_PRIVATE_KEY_PEM) for an API account that can read audit events.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/auditevents"
)

func main() {
	c, err := axm.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	now := time.Now().UTC()
	response, _, err := c.AXMAPI.AuditEvents.GetV1(ctx, &auditevents.RequestQueryOptions{
		FilterStartTimestamp: now.Add(-24 * time.Hour).Format(time.RFC3339),
		FilterEndTimestamp:   now.Format(time.RFC3339),
		FilterType:           auditevents.AuditEventTypeTermsAccepted,
		Fields:               []string{auditevents.FieldType, auditevents.FieldEventDataPropertyKey, auditevents.FieldEventDataTermsAccepted},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d terms acceptance events\n", len(response.Data))
	for _, event := range response.Data {
		if attrs := event.Attributes; attrs != nil && attrs.EventDataTermsAccepted != nil {
			fmt.Printf("Accepted terms: %s\n", attrs.EventDataTermsAccepted.TermsName)
		}
	}
}
