package pbkdf2

import (
	p "github.com/pganguli/pbkdf2"
)

var params = &p.Params{
	Iterations:  210000,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	hash, err := p.CreateHash(password, params)
	return hash, err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password string, hash string) (bool, error) {
	match, err := p.ComparePasswordAndHash(password, hash)
	return match, err
}
