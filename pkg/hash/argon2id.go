package hash

import (
	"github.com/alexedwards/argon2id"
)

var params = &argon2id.Params{
	Memory:      64 * 1024, // m
	Iterations:  1,         // t
	Parallelism: 2,         // p
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, params)
	return hash, err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match, err
}
