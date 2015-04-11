// Package signed generates signed UUIDs.
// It generates V4 UUIDs and signs them
// with a secret using a SHA256 HMAC.
package signed

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/twinj/uuid"
)

// Signed is a container for an UUID and a signature
type Signed struct {
	id   uuid.UUID
	sign []byte
}

func newMac(src []byte, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(src)
	return mac.Sum(nil)
}

// New creates a new UUID and signs it using the given secret.
func New(secret []byte) *Signed {
	id := uuid.NewV4()
	sign := newMac(id.Bytes(), secret)
	return &Signed{id, sign}
}

// Parse parses a string and returns a signed UUID
// or an error if the format is invalid.
// Note that the UUID is not being validated.
func Parse(src string) (*Signed, error) {
	x := strings.Split(src, ":")

	if len(x) != 2 {
		return nil, errors.New("invalid format")
	}

	id, err := uuid.ParseUUID(x[0])
	if err != nil {
		return nil, err
	}

	sign, err := base64.StdEncoding.DecodeString(x[1])
	if err != nil {
		return nil, err
	}

	return &Signed{id, sign}, nil
}

// ID returns the ID part of a signed UUID
func (s Signed) ID() string {
	return uuid.Formatter(s.id, uuid.CleanHyphen)
}

// Signature returns the signature of a signed UUID
func (s Signed) Signature() string {
	return base64.StdEncoding.EncodeToString(s.sign)
}

func (s Signed) String() string {
	return s.ID() + ":" + s.Signature()
}

// Validate validates the signed UUID against a given secret.
func (s Signed) Validate(secret []byte) bool {
	expected := newMac(s.id.Bytes(), secret)
	return hmac.Equal(s.sign, expected)
}
