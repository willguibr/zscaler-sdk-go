package integration

import (
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/zia/services/adminuserrolemgmt"
)

func TestAdminRoles_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := adminuserrolemgmt.New(client)

	roles, err := service.GetAllAdminRoles()
	if err != nil {
		t.Errorf("Error getting admin roles: %v", err)
		return
	}
	if len(roles) == 0 {
		t.Errorf("No admin roles found")
		return
	}
	name := roles[0].Name
	t.Log("Getting admin roles by name:" + name)
	role, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting admin roles by name: %v", err)
		return
	}
	if role.Name != name {
		t.Errorf("admin role name does not match: expected %s, got %s", name, role.Name)
		return
	}
}
