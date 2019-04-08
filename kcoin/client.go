package kcoin

import "github.com/dghubble/sling"

// GetClient returns a kalicoin http client
func GetClient() *sling.Sling {
	return sling.New().
		Base(config.URI).
		SetBasicAuth(config.Username, config.Password)
}
