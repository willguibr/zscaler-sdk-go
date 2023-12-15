package dlp_exact_data_match

import (
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v2/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPEDM_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	templates, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting idm profiles: %v", err)
		return
	}
	if len(templates) == 0 {
		t.Errorf("No idm profile found")
		return
	}
	name := templates[0].ProjectName
	t.Log("Getting edm template by name:" + name)
	template, err := service.GetDLPEDMByName(name)
	if err != nil {
		t.Errorf("Error getting edm template by name: %v", err)
		return
	}
	if template.ProjectName != name {
		t.Errorf("edm template name does not match: expected %s, got %s", name, template.ProjectName)
		return
	}
	// Negative Test: Try to retrieve an edm template with a non-existent name
	nonExistentName := "ThisEDMTemplateDoesNotExist"
	_, err = service.GetDLPEDMByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestGetById(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Get all servers to find a valid ID
	templates, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error getting all edm templates: %v", err)
	}
	if len(templates) == 0 {
		t.Fatalf("No edm templates found for testing")
	}

	// Choose the first server's ID for testing
	testID := templates[0].SchemaID

	// Retrieve the server by ID
	template, err := service.GetDLPEDMSchemaID(testID)
	if err != nil {
		t.Errorf("Error retrieving edm template with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved server
	if template == nil {
		t.Errorf("No server returned for ID %d", testID)
		return
	}

	if template.SchemaID != testID {
		t.Errorf("Retrieved server ID mismatch: expected %d, got %d", testID, template.SchemaID)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	templates, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting edm template: %v", err)
		return
	}
	if len(templates) == 0 {
		t.Errorf("No edm template found")
		return
	}

	// Validate edm template
	for _, template := range templates {
		// Checking if essential fields are not empty
		if template.SchemaID == 0 {
			t.Errorf("edm template ID is empty")
		}
		if template.ProjectName == "" {
			t.Errorf("edm template Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming a edm template with the name "BD_EDM_TEMPLATE01" exists
	knownName := "BD_EDM_TEMPLATE01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		template, err := service.GetDLPEDMByName(variation)
		if err != nil {
			t.Errorf("Error getting edm template with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if template.ProjectName != knownName {
			t.Errorf("Expected edm template name to be '%s' for variation '%s', but got '%s'", knownName, variation, template.ProjectName)
		}
	}
}

func TestEDMFields(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Retrieve all EDM profiles
	edmProfiles, err := service.GetAll() // Assuming appropriate method name and parameters
	if err != nil {
		t.Fatalf("Error getting all EDM profiles: %v", err)
	}
	if len(edmProfiles) == 0 {
		t.Fatalf("No EDM profiles found for testing")
	}

	// Iterate through each EDM profile and check various fields
	for _, profile := range edmProfiles {
		if profile.SchemaID == 0 {
			t.Errorf("SchemaID field is empty")
		}
		if profile.ProjectName == "" {
			t.Errorf("ProjectName field is empty")
		}
		if profile.SchemaStatus != "EDM_INDEXING_SUCCESS" {
			t.Errorf("SchemaStatus field is not as expected: got %s", profile.SchemaStatus)
		}
		if !profile.SchemaActive {
			t.Errorf("SchemaActive field is not active")
		}

		// Asserting elements in the TokenList
		for _, token := range profile.TokenList {
			if token.Name == "" || token.Type == "" {
				t.Errorf("Token fields Name or Type are not properly populated")
			}
			// Here, we check if at least one token is marked as primary key
			if token.PrimaryKey {
				t.Log("A primary key token found")
				break // Break after finding the first primary key
			}
		}
	}
}
