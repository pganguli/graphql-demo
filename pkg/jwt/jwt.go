package jwt

import (
	"crypto"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"os"
	"time"
)

func loadData(path string) ([]byte, error) {
	var reader io.Reader
	if file, err := os.Open(path); err == nil {
		reader = file
		defer file.Close()
	} else {
		return nil, err
	}

	return io.ReadAll(reader)
}

type KeyPair struct {
	privKey *crypto.PrivateKey
	pubKey  *crypto.PublicKey
}

func KeyFromPEM(keyPair *KeyPair) error {
	if data, err := loadData(os.Getenv("PRIVATE_KEY")); err == nil {
		if key, err := jwt.ParseEdPrivateKeyFromPEM(data); err == nil {
			keyPair.privKey = &key
		} else {
			return err
		}
	} else {
		return err
	}

	if data, err := loadData(os.Getenv("PUBLIC_KEY")); err == nil {
		if key, err := jwt.ParseEdPublicKeyFromPEM(data); err == nil {
			keyPair.pubKey = &key
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var keyPair = new(KeyPair)

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	if keyPair.privKey == nil {
		if err := KeyFromPEM(keyPair); err != nil {
			return "", fmt.Errorf("Couldn't read key: %w", err)
		}
	}

	tokenString, err := token.SignedString(*keyPair.privKey)
	if err != nil {
		log.Fatal("Error in signing token")
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username it it's claims
func ParseToken(tokenString string) (string, error) {
	if keyPair.pubKey == nil {
		if err := KeyFromPEM(keyPair); err != nil {
			return "", fmt.Errorf("Couldn't read key: %w", err)
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return *keyPair.pubKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.Username, nil
	} else {
		return "", err
	}
}
