package kcoin

import "github.com/dghubble/sling"

func GetClient() *sling.Sling {
	return sling.New().
		Base(config.URI).
		SetBasicAuth(config.Username, config.Password)
}
