package rbac

import (
	"github.com/casbin/casbin/v2"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
)

type casbinManager struct {
	Enforcer *casbin.Enforcer
}

var CasbinManager casbinManager

func Init(_ef *casbin.Enforcer) {
	CasbinManager = casbinManager{Enforcer: _ef}
}

func (c *casbinManager) SetUserUnderGroup(userID, groupID string, action string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)
	_, err := c.Enforcer.AddNamedPolicy("p", groupString, userString, action)
	if err != nil {
		return err
	}
	return nil
}

func (c *casbinManager) SetUserToGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)

	_, err := c.Enforcer.AddGroupingPolicy(userString, groupString)
	if err != nil {
		return err
	}

	return nil
}

func (c *casbinManager) RemoveUserFromGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)

	_, err := c.Enforcer.RemoveGroupingPolicy(userString, groupString)
	if err != nil {
		return err
	}

	return nil
}

func (c *casbinManager) SetUserToReadGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)

	_, err := c.Enforcer.AddNamedPolicy("p", userString, groupString, rbacValues.Read)
	if err != nil {
		return err
	}

	return nil
}

func (c *casbinManager) SetUserToWriteGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)

	_, err := c.Enforcer.AddNamedPolicy("p", userString, groupString, rbacValues.Write)
	if err != nil {
		return err
	}

	return nil
}

func (c *casbinManager) CheckUserToUserWritePermission(userID, userID2 string) (bool, error) {
	userString := rbacValues.UserIDPrefix(userID)
	userString2 := rbacValues.UserIDPrefix(userID2)
	return c.Enforcer.Enforce(userString, userString2, rbacValues.Write)
}

func (c *casbinManager) CheckUserToUserReadPermission(userID, userID2 string) (bool, error) {
	userString := rbacValues.UserIDPrefix(userID)
	userString2 := rbacValues.UserIDPrefix(userID2)
	return c.Enforcer.Enforce(userString, userString2, rbacValues.Read)
}

func (c *casbinManager) SetUserWritePermissionToGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)
	_, err := c.Enforcer.AddNamedPolicy("p", userString, groupString, rbacValues.Write)
	if err != nil {
		return err
	}
	return nil
}

func (c *casbinManager) CheckUserWritePermissionToGroup(userID, groupID string) (bool, error) {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)
	return c.Enforcer.Enforce(userString, groupString, rbacValues.Write)
}

func (c *casbinManager) SetUserReadPermissionToGroup(userID, groupID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)
	_, err := c.Enforcer.AddNamedPolicy("p", userString, groupString, rbacValues.Read)
	if err != nil {
		return err
	}
	return nil
}

func (c *casbinManager) CheckUserReadPermissionToGroup(userID, groupID string) (bool, error) {
	userString := rbacValues.UserIDPrefix(userID)
	groupString := rbacValues.GroupIDPrefix(groupID)
	return c.Enforcer.Enforce(userString, groupString, rbacValues.Read)
}

// func (c *casbinManager) CheckUserReadPermissionToUserInGroup(userID, groupID string) (bool, error) {
// 	userString := rbacValues.UserIDPrefix(userID)
// 	groupString := rbacValues.GroupIDPrefix(groupID)
// 	return c.Enforcer.Enforce(userString, groupString, rbacValues.Read)
// }

func (c *casbinManager) SetSubgroup(parentGroup, subGroup string) error {
	parentGroupString := rbacValues.GroupIDPrefix(parentGroup)
	subGroupString := rbacValues.GroupIDPrefix(subGroup)
	_, err := c.Enforcer.AddGroupingPolicy(subGroupString, parentGroupString)
	if err != nil {
		return err
	}
	return nil
}

// SetUserGlobalAdmin Global Admin Rights
func (c *casbinManager) SetUserGlobalAdmin(userID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	_, err := c.Enforcer.AddGroupingPolicy(userString, rbacValues.ADMIN_PREFIX)
	return err
}

func (c *casbinManager) RemoveUserGlobalAdmin(userID string) error {
	userString := rbacValues.UserIDPrefix(userID)
	_, err := c.Enforcer.RemoveGroupingPolicy(userString, rbacValues.ADMIN_PREFIX)
	return err
}

func (c *casbinManager) CheckUserIsGlobalAdmin(userID string) (bool, error) {
	// Let's enforce if the user has wildcard "*" access on resource "*", action "invoke"
	// which we can define as the "Admin" role policy
	userString := rbacValues.UserIDPrefix(userID)
	return c.Enforcer.Enforce(userString, "*", "invoke")
}

// SetupGlobalAdminPolicy seeds the DB with the wildcard policy for admins
func (c *casbinManager) SetupGlobalAdminPolicy() error {
	_, err := c.Enforcer.AddNamedPolicy("p", rbacValues.ADMIN_PREFIX, "*", "invoke")
	return err
}

// CheckUserPermission is a quick helper for manually evaluating permissions
// commonly called inside gRPC services.
func (c *casbinManager) CheckUserPermission(userID string, resource string, action string) (bool, error) {
	userString := rbacValues.UserIDPrefix(userID)
	return c.Enforcer.Enforce(userString, resource, action)
}
