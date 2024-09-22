package policysetcontrollerv2

import (
	"fmt"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccessInspectionPolicyInspectV2(t *testing.T) {
	policyType := "INSPECTION_POLICY"
	inspectionProfileID := "BD_SA_Profile1"
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	idpList, _, err := idpcontroller.GetAll(service)
	if err != nil {
		t.Errorf("Error getting idps: %v", err)
		return
	}
	if len(idpList) == 0 {
		t.Skip("No IdPs retrieved, skipping test as it requires at least one IdP")
		return
	}

	samlsList, _, err := samlattribute.GetAll(service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(samlsList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	profileID, _, err := inspection_profile.GetByName(service, inspectionProfileID)
	if err != nil {
		t.Errorf("Error getting inspection profile id set: %v", err)
		return
	}

	accessPolicySet, _, err := GetByPolicyType(service, policyType)
	if err != nil {
		t.Errorf("Error getting access policy set: %v", err)
		return
	}

	var ruleIDs []string // Store the IDs of the created rules

	for i := 0; i < 5; i++ {
		// Generate a unique name for each iteration
		name := fmt.Sprintf("tests-%s-%d", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha), i)

		accessPolicyRule := PolicyRule{
			Name:                   name,
			Description:            name,
			PolicySetID:            accessPolicySet.ID,
			Action:                 "INSPECT",
			ZpnInspectionProfileID: profileID.ID,
			Conditions: []PolicyRuleResourceConditions{
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "SAML",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: samlsList[0].ID, RHS: "user1@acme.com"}},
						},
					},
				},
				{
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "PLATFORM",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: "linux", RHS: "true"}, {LHS: "windows", RHS: "true"}},
						},
					},
				},
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_exporter"},
						},
					},
				},
			},
		}

		// Test resource creation
		createdResource, _, err := CreateRule(service, &accessPolicyRule)

		if err != nil {
			t.Errorf("Error making POST request: %v", err)
			continue
		}
		if createdResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
			continue
		}

		ruleIDs = append(ruleIDs, createdResource.ID) // Collect rule ID for reordering

		// Update the rule name
		updatedName := name + "-updated"
		accessPolicyRule.Name = updatedName
		_, updateErr := UpdateRule(service, accessPolicySet.ID, createdResource.ID, &accessPolicyRule)

		if updateErr != nil {
			t.Errorf("Error updating rule: %v", updateErr)
			continue
		}

		// Retrieve and print the updated resource as JSON
		updatedResource, _, getErr := GetPolicyRule(service, accessPolicySet.ID, createdResource.ID)
		if getErr != nil {
			t.Errorf("Error retrieving updated resource: %v", getErr)
			continue
		}
		if updatedResource.Name != updatedName {
			t.Errorf("Expected updated resource name '%s', but got '%s'", updatedName, updatedResource.Name)
		}

		// Test resource retrieval by name
		updatedResource, _, err = GetByNameAndType(service, policyType, updatedName)
		if err != nil {
			t.Errorf("Error retrieving resource by name: %v", err)
		}
		if updatedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
		}
		if updatedResource.Name != updatedName {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", updatedName, updatedResource.Name)
		}
		time.Sleep(5 * time.Second)
	}

	// Reorder the rules after all have been created and updated
	ruleIdToOrder := make(map[string]int)
	for i, id := range ruleIDs {
		ruleIdToOrder[id] = len(ruleIDs) - i // Reverse the order
	}

	_, err = BulkReorder(service, policyType, ruleIdToOrder)
	if err != nil {
		t.Errorf("Error reordering rules: %v", err)
	}

	// Optionally verify the new order of rules here

	// Clean up: Delete the rules
	for _, ruleID := range ruleIDs {
		_, err = Delete(service, accessPolicySet.ID, ruleID)
		if err != nil {
			t.Errorf("Error deleting resource: %v", err)
		}
	}
}