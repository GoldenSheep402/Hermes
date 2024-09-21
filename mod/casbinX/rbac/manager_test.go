// manager_test.go
package rbac

import (
	"testing"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func getTestModel() string {
	return rbacValues.RbacRule
}

func initEnforcer() (*casbin.Enforcer, error) {
	m, err := model.NewModelFromString(getTestModel())
	if err != nil {
		return nil, err
	}

	// Use a memory adapter for policies
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func TestCasbinManager(t *testing.T) {
	enforcer, err := initEnforcer()
	if err != nil {
		t.Fatalf("Failed to initialize enforcer: %v", err)
	}

	// Initialize the casbinManager
	Init(enforcer)

	// Test SetUserToGroup
	if err := CasbinManager.SetUserToGroup("user1", "group1"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}
	if err := CasbinManager.SetUserToGroup("user2", "group1"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}

	// Verify grouping policy
	if ok, _ := enforcer.HasGroupingPolicy(rbacValues.UserIDPrefix("user1"), rbacValues.GroupIDPrefix("group1")); !ok {
		t.Errorf("Grouping policy not added")
	}
	if err := CasbinManager.SetUserToGroup("user1", "group1"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}

	// Test RemoveUserFromGroup
	if err := CasbinManager.RemoveUserFromGroup("user1", "group1"); err != nil {
		t.Errorf("RemoveUserFromGroup failed: %v", err)
	}

	// Verify grouping policy removed
	if ok, _ := enforcer.HasGroupingPolicy(rbacValues.UserIDPrefix("user1"), rbacValues.GroupIDPrefix("group1")); ok {
		t.Errorf("Grouping policy not removed")
	}

	// Test SetUserUnderGroup
	if err := CasbinManager.SetUserUnderGroup("user2", "group1", rbacValues.Read); err != nil {
		t.Errorf("SetUserUnderGroup failed: %v", err)
	}

	// Verify policy added
	if ok, _ := enforcer.HasPolicy(rbacValues.GroupIDPrefix("group1"), rbacValues.UserIDPrefix("user2"), rbacValues.Read); !ok {
		t.Errorf("Policy not added")
	}

	if err := CasbinManager.SetUserToGroup("user4", "group1"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}
	// Test CheckUserToUserReadPermission
	allowed, err := CasbinManager.CheckUserToUserReadPermission("user4", "user2")
	if err != nil {
		t.Errorf("CheckUserToUserReadPermission failed: %v", err)
	}
	if !allowed {
		t.Errorf("Expected permission to be allowed")
	}

	// Test CheckUserToUserWritePermission (should be denied)
	allowed, err = CasbinManager.CheckUserToUserWritePermission("user1", "group1")
	if err != nil {
		t.Errorf("CheckUserToUserWritePermission failed: %v", err)
	}
	if allowed {
		t.Errorf("Expected permission to be denied")
	}

	// Add write permission and test again
	if err := CasbinManager.SetUserUnderGroup("user3", "group1", rbacValues.Write); err != nil {
		t.Errorf("SetUserUnderGroup failed: %v", err)
	}

	if err := CasbinManager.SetUserToGroup("user5", "group1"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}
	allowed, err = CasbinManager.CheckUserToUserWritePermission("user5", "user3")
	if err != nil {
		t.Errorf("CheckUserToUserWritePermission failed: %v", err)
	}
	if !allowed {
		t.Errorf("Expected permission to be allowed")
	}

	// Test SetSubgroup
	if err := CasbinManager.SetSubgroup("group1", "group2"); err != nil {
		t.Errorf("SetSubgroup failed: %v", err)
	}

	// Verify subgroup policy
	if ok, _ := enforcer.HasGroupingPolicy(rbacValues.GroupIDPrefix("group2"), rbacValues.GroupIDPrefix("group1")); !ok {
		t.Errorf("Subgroup policy not added")
	}

	// Add user to subgroup
	if err := CasbinManager.SetUserToGroup("user2", "group2"); err != nil {
		t.Errorf("SetUserToGroup failed: %v", err)
	}

	// // Test CheckUserReadPermissionToUserInGroup
	// allowed, err = CasbinManager.CheckUserReadPermissionToUserInGroup("user2", "group1")
	// if err != nil {
	// 	t.Errorf("CheckUserReadPermissionToUserInGroup failed: %v", err)
	// }
	// if !allowed {
	// 	t.Errorf("Expected permission to be allowed through subgroup")
	// }
}
