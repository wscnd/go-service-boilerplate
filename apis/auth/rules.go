package auth

import _ "embed"

// Package name of our rego code.
const (
	opaPackage string = "serviceAuth.rego"
)

// These the current set of rules we have for auth.
const (
	RuleAuthenticate   = "auth"
	RuleAny            = "rule_any"
	RuleAdminOnly      = "rule_admin_only"
	RuleUserOnly       = "rule_user_only"
	RuleAdminOrSubject = "rule_admin_or_subject"
)

// Core OPA policies.
var (
	//go:embed rego/authentication.rego
	opaAuthentication string

	//go:embed rego/authorization.rego
	opaAuthorization string
)
