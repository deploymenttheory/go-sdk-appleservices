// Set APPLE_KEY_ID, APPLE_ISSUER_ID and APPLE_PRIVATE_KEY_PATH (or
// APPLE_PRIVATE_KEY_PEM), then pass the device ID as the only argument.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: GetActivationLockStatusByDeviceID <device-id>")
	}
	c, err := axm.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, _, err := c.AXMAPI.Devices.GetActivationLockStatusByDeviceIDV1(ctx, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	attrs := response.Data.Attributes
	if attrs == nil || attrs.IsLocked == nil {
		fmt.Println("Activation Lock status was not supplied")
		return
	}
	fmt.Printf("Locked: %t; type: %s\n", *attrs.IsLocked, attrs.LockType)
}
