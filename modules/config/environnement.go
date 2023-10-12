package config

// Environnement is an enum for environnement
// It can be "development", "production" or "test"
type Environnement string

const (
	// Production environnement
	Production Environnement = "production"
	// Development environnement
	Development Environnement = "development"
	// Staging environnement
	Staging Environnement = "staging"
)

// String returns the string representation of the environnement
func (e Environnement) String() string {
	return string(e)
}

// IsProduction returns true if the environnement is production
func (e Environnement) IsProduction() bool {
	return e == Production
}

// IsDevelopment returns true if the environnement is development
func (e Environnement) IsDevelopment() bool {
	return e == Development
}

// IsStaging returns true if the environnement is staging
func (e Environnement) IsStaging() bool {
	return e == Staging
}

// IsLocal returns true if the environnement is development or test
func (e Environnement) IsLocal() bool {
	return e.IsDevelopment() || e.IsStaging()
}

// IsRemote returns true if the environnement is production
func (e Environnement) IsRemote() bool {
	return e.IsProduction()
}
