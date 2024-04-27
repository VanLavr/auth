package hasher

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	salt int
}

var Hshr = NewHasher(4)

func NewHasher(salt int) *Hasher {
	return &Hasher{salt: salt}
}

func (h *Hasher) Hash(data string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), h.salt)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

func (h *Hasher) Validate(origin string, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(origin), []byte(data))
	return err == nil
}
