package utils

import (
	"github.com/alexedwards/argon2id"
)

func GenerateHash(password string) (string, error) {
	return argon2id.CreateHash("pa$$word", argon2id.DefaultParams)

}

func CompareHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
