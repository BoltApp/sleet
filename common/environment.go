package common

// Environment provides a common way of interacting with Sleet's PsP
// Sandbox refers to non-live, typically test accounts and Production to live accounts
// Done at the Sleet level to avoid clients having to import Payment specific data
type Environment string

// Only sandbox and production environments are supported
const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)
