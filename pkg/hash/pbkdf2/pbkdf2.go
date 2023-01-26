package pbkdf2

import (
	p "github.com/pganguli/pbkdf2"
)

var params = &p.Params{
	Iterations:  210000,
	SaltLength:  16,
	KeyLength:   64,
}

// HashPassword hashes given password
func CreateHash(password string) (string, error) {
	return p.CreateHash(password, params)
}

// CheckPassword hash compares raw password with it's hashed values
func ComparePasswordAndHash(password string, hash string) (bool, error) {
	return p.ComparePasswordAndHash(password, hash)
}
