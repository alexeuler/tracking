package models

import "github.com/up-finder/silk.web/app/db"

const (
	AUTH_HASH_NAME = "silk:auth:keys" //hash in redis where domain-secret pairs are stored
)

type Auth struct {
	Domain string
	Key    string //secret
}

// Get Auth data with specified domain
func FetchAuth(domain string) *Auth {
	key, _ := db.Redis.HGet(AUTH_HASH_NAME, domain).Result()
	return &Auth{Domain: domain, Key: key}
}
