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
	ADMIN_PREFIX = "Admin"
	USER_PREFIX  = "User"
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
	Level8
	Level9
	Level10
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

// Regex pattern shared by SplitID and SplitIdAndLevel
var sharedPattern = regexp.MustCompile(`^([^:]+):([^:]+)(?::([^:]+))?$`)

// SplitID splits the permission into object and id.
// For example, given "User:123", it returns "User", "123", nil.
func SplitID(permission string) (string, string, error) {
	matches := sharedPattern.FindStringSubmatch(permission)
	if len(matches) < 3 {
		return "", "", errors.New("permission format is invalid")
	}
	object := matches[1]
	id := matches[2]
	return object, id, nil
}

// SplitIdAndLevel splits the role permission into prefix, id, and level.
// For example, given "Group:456:2", it returns "Group", "456", "2", nil.
// If the level is not provided, level will be an empty string.
func SplitIdAndLevel(rolePermission string) (string, string, string, error) {
	matches := sharedPattern.FindStringSubmatch(rolePermission)
	if len(matches) < 3 {
		return "", "", "", errors.New("role permission format is invalid")
	}
	prefix := matches[1]
	id := matches[2]
	level := ""
	if len(matches) >= 4 {
		level = matches[3]
	}
	return prefix, id, level, nil
}
