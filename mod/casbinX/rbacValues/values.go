package rbacValues

import (
	"errors"
	"regexp"
	"strconv"
)

const RbacRule = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || (r.act == "read" && p.act == "write"))
`

const (
	Read  = "read"
	Write = "write"
)

// Prefixes
const (
	ADMIN_PREFIX = "ADMIN"
	USER_PREFIX  = "USER"
)

// Levels
const (
	Level0 = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
)

// ID Prefix Functions
func UserIDPrefix(userID string) string {
	return USER_PREFIX + ":" + userID
}

func GroupIDPrefix(groupID string) string {
	return "Group:" + groupID
}

func GroupWithIDAndLevelPrefix(groupID string, level int) string {
	return "Group:" + groupID + ":" + strconv.Itoa(level)
}

func CategoryIDPrefix(categoryID string) string {
	return "Category:" + categoryID
}

func CategoryWithIDAndLevelPrefix(categoryID string, level int) string {
	return "Category:" + categoryID + ":" + strconv.Itoa(level)
}

// Regex pattern shared by SplitPermission and SplitIdAndLevel
var sharedPattern = regexp.MustCompile(`^([^:]+):([^:]+)(?::([^:]+))?`)

// SplitID splits the permission into object and id.
func SplitID(permission string) (string, string, error) {
	matches := sharedPattern.FindStringSubmatch(permission)
	if len(matches) != 3 && len(matches) != 4 { // Handles both cases
		return "", "", errors.New("permission format is invalid")
	}
	return matches[1], matches[2], nil
}

// SplitIdAndLevel splits the role permission into prefix, id, and optionally level.
func SplitIdAndLevel(rolePermission string) (string, string, string, error) {
	matches := sharedPattern.FindStringSubmatch(rolePermission)
	if len(matches) != 4 {
		return "", "", "", errors.New("role permission format is invalid")
	}
	return matches[1], matches[2], matches[3], nil
}
