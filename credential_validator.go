package main

// We can use this interface to find & validate credentials of any
// other providers
type CredentialValidator interface {
	FindCredentials(content string, pattern string)
	ValidateCredentials(key string, token string)
}
