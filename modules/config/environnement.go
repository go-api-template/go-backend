package config

// Environnement is an enum for environnement
// It can be "development", "production" or "test"
type Environnement string

const (
	// Development environnement
	Development Environnement = "development"
	// Production environnement
	Production Environnement = "production"
	// Test environnement
	Test Environnement = "test"
)

// String returns the string representation of the environnement
func (e Environnement) String() string {
	return string(e)
}

//// ParseRole parses a string into a Role
//func ParseRole(s string) (r Role, err error) {
//	roles := map[Role]struct{}{
//		RoleAdmin: {},
//		RoleUser:  {},
//	}
//
//	r = Role(s)
//	_, ok := roles[r]
//	if !ok {
//		return r, fmt.Errorf(`cannot parse:[%s] as Role`, s)
//	}
//	return r, nil
//}

// IsDevelopment returns true if the environnement is development
func (e Environnement) IsDevelopment() bool {
	return e == Development
}

// IsProduction returns true if the environnement is production
func (e Environnement) IsProduction() bool {
	return e == Production
}

// IsTest returns true if the environnement is test
func (e Environnement) IsTest() bool {
	return e == Test
}

// IsLocal returns true if the environnement is development or test
func (e Environnement) IsLocal() bool {
	return e.IsDevelopment() || e.IsTest()
}

// IsRemote returns true if the environnement is production
func (e Environnement) IsRemote() bool {
	return e.IsProduction()
}
