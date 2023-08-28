package policysetcontroller

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/tests"
)

var supportedPolicyTypes = []string{
	"ACCESS_POLICY",
	"TIMEOUT_POLICY",
	"CLIENT_FORWARDING_POLICY",
	"ISOLATION_POLICY",
}

// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	if !present {
		return true // default value
	}
	shouldClean, err := strconv.ParseBool(val)
	if err != nil {
		return true // default to cleaning if the value is not parseable
	}
	log.Printf("ZSCALER_SDK_TEST_SWEEP value: %v", shouldClean)
	return shouldClean
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	for _, policyType := range supportedPolicyTypes {
		resources, _, err := service.GetAllByType(policyType)
		if err != nil {
			log.Printf("Error fetching resources of type %s: %v", policyType, err)
			continue
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %s, Name: %s, Type: %s", r.ID, r.Name, policyType)

			// Assuming that the Delete function needs both policySetID and policyRuleID
			_, _ = service.Delete(r.PolicySetID, r.ID)
		}
	}
}