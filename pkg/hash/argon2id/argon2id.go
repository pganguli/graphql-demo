package argon2id

import (
	a "github.com/alexedwards/argon2id"
)

var params = &a.Params{
	Memory:      64 * 1024, // m
	Iterations:  1,         // t
	Parallelism: 2,         // p
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword hashes given password
func CreateHash(password string) (string, error) {
	return a.CreateHash(password, params)
}

// CheckPassword hash compares raw password with it's hashed values
func ComparePasswordAndHash(password string, hash string) (bool, error) {
	return a.ComparePasswordAndHash(password, hash)
}
