package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{Cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *BcryptHasher) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type Hasher interface {
	Hash(password string) (string, error)
	CheckPasswordHash(password, hash string) error
}

type HMACHasher struct {
	Secret []byte
}

func NewHMACHasher(secret string) *HMACHasher {
	return &HMACHasher{
		Secret: []byte(secret),
	}
}

func (h *HMACHasher) Hash(value string) string {
	mac := hmac.New(sha256.New, h.Secret)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

func (h *HMACHasher) Compare(value, hash string) bool {
	return h.Hash(value) == hash
}

type CodeHasher interface {
	Hash(value string) string
	Compare(value, hash string) bool
}
