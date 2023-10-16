package models

import "fmt"

// Role given to a user
type Role string

// Enum of types of roles
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// ParseRole parses a string into a Role
func ParseRole(s string) (r Role, err error) {
	roles := map[Role]struct{}{
		RoleAdmin: {},
		RoleUser:  {},
	}

	r = Role(s)
	_, ok := roles[r]
	if !ok {
		return r, fmt.Errorf(`cannot parse:[%s] as Role`, s)
	}
	return r, nil
}

// IsAdmin checks if the role is admin
func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

// IsUser checks if the role is user
func (r Role) IsUser() bool {
	return r == RoleUser
}
