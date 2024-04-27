package hasher

import (
	"crypto/sha512"
	"hash"
)

type Hasher struct {
	hasher hash.Hash
}

var Hshr = NewHasher()

func NewHasher() *Hasher {
	return &Hasher{hasher: sha512.New()}
}

func (h *Hasher) Encrypt(data string) string {
	defer h.hasher.Reset()
	h.hasher.Write([]byte(data))

	return string(h.hasher.Sum(nil))
}

func (h *Hasher) Validate(origin string, data string) bool {
	defer h.hasher.Reset()

	h.hasher.Write([]byte(data))

	return string(h.hasher.Sum(nil)) == origin
}
