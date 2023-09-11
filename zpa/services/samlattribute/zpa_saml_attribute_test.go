package samlattribute

import (
	"testing"

	tests "github.com/SecurityGeekIO/zscaler-sdk-go/v2/tests"
)

func TestSAMLAttribute(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	attributes, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(attributes) == 0 {
		t.Errorf("No saml attribute found")
		return
	}
	name := attributes[0].Name
	t.Log("Getting saml attribute by name:" + name)
	attribute, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting saml attribute by name: %v", err)
		return
	}
	if attribute.Name != name {
		t.Errorf("identity provider name does not match: expected %s, got %s", name, attribute.Name)
		return
	}
}
