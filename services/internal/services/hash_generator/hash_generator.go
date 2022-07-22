package hashgenerator

import (
	"go_rest/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type HashGeneratorInterface interface {
	Hash(pass string) (string, error)
	Check(pass, hash string) bool
}

type HashGenerator struct {
	salt string
}

func (h *HashGenerator) Salt(pass string) string {
	// TODO сложно сказать нужно ли солить bcrypt
	return pass + h.salt
}

func (h *HashGenerator) Hash(pass string) (string, error) {
	saltedPass := h.Salt(pass)
	res, err := bcrypt.GenerateFromPassword([]byte(saltedPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (h *HashGenerator) Check(pass, hash string) bool {
	saltedPass := h.Salt(pass)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPass))
	return err == nil
}

func NewHashGenerator(c *config.Config) *HashGenerator {
	return &HashGenerator{
		salt: c.Salt,
	}
}
