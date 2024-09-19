package manager

import (
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/casbin/casbin/v2"
)

type casbinManager struct {
	Enforcer *casbin.Enforcer
}

var CasbinManager casbinManager

func Init(_ef *casbin.Enforcer) {
	CasbinManager = casbinManager{Enforcer: _ef}
}

// InitPolicy initializes the policy of the user.
// The policy is based on the user's role and the project or product.
func (c *casbinManager) InitPolicy(userID string, prefixedSub string) error {
	_, err := c.Enforcer.AddNamedPolicy("p", prefixedSub+":ADMIN", prefixedSub, rbacValues.Write)
	if err != nil {
		return err
	}
	_, err = c.Enforcer.AddNamedPolicy("p", prefixedSub+":MEMBER", prefixedSub, rbacValues.Read)
	if err != nil {
		return err
	}
	_, err = c.Enforcer.AddGroupingPolicy(rbacValues.UserIDPrefix(userID), prefixedSub+":ADMIN")
	if err != nil {
		return err
	}
	return nil
}

func (c *casbinManager) CheckProjectWritePermission(userID, projectID string) (bool, error) {
	// return c.Enforcer.Enforce(rbacValues.UserIDPrefix(userID), rbacValues.ProjectIDPrefix(projectID), rbacValues.Write)
	return true, nil
}

func (c *casbinManager) CheckProjectReadPermission(userID, projectID string) (bool, error) {
	// return c.Enforcer.Enforce(rbacValues.UserIDPrefix(userID), rbacValues.ProjectIDPrefix(projectID), rbacValues.Read)
	return true, nil
}
