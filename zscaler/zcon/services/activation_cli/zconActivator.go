package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon/services/activation"
)

func main() {
	log.Printf("[INFO] Initializing ZCON client\n")

	// Attempt to initialize the client only at runtime
	err := runZCONActivator()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

// Separate runtime logic
func runZCONActivator() error {
	username := os.Getenv("ZCON_USERNAME")
	password := os.Getenv("ZCON_PASSWORD")
	apiKey := os.Getenv("ZCON_API_KEY")
	zconCloud := os.Getenv("ZCON_CLOUD")

	// Validate credentials at runtime
	if username == "" || password == "" || apiKey == "" || zconCloud == "" {
		return fmt.Errorf("missing ZCON credentials: ensure ZCON_USERNAME, ZCON_PASSWORD, ZCON_API_KEY and ZCON_CLOUD environment variables are set")
	}

	// Create a new ZCON configuration
	zconCfg, err := zcon.NewConfiguration(
		zcon.WithZconUsername(username),
		zcon.WithZconPassword(password),
		zcon.WithZconAPIKey(apiKey),
		zcon.WithZconCloud(zconCloud),
		zcon.WithUserAgentExtra(""),
	)
	if err != nil {
		return fmt.Errorf("failed to create ZCON configuration: %w", err)
	}

	// Initialize ZCON client
	zconClient, err := zcon.NewClient(zconCfg)
	if err != nil {
		return fmt.Errorf("failed to create ZCON client: %w", err)
	}

	// Wrap the ZCON client in a Service instance
	service := services.New(zconClient)

	resp, err := activation.ForceActivationStatus(service, activation.ECAdminActivation{
		OrgEditStatus:         "org_edit_status",
		OrgLastActivateStatus: "org_last_activate_status",
		AdminActivateStatus:   "admin_activate_status",
	})
	if err != nil {
		log.Printf("[ERROR] Activation Failed: %v\n", err)
	} else {
		log.Printf("[INFO] Activation succeeded: %#v\n", resp)
	}
	log.Printf("[INFO] Destroying session: %#v\n", resp)

	// Log out the client
	zconClient.Logout()

	return nil
}
