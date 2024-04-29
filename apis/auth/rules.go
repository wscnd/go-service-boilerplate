package auth

import _ "embed"

// Package name of our rego code.
const (
	OpaPackage string = "serviceAuth.rego"
)

// These the current set of rules we have for auth.
const (
	RuleAuthenticate = "auth"
)

// Core OPA policies.
var (
	//go:embed rego/authentication.rego
	OpaAuthentication string
)
