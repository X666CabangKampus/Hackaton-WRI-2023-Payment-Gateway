package util

import (
	"crypto/sha256"
	"fmt"
)

const SECRET = "hXwGc-Niw.hq,fu:v+HK8|_uVA;xunkGB$u5]NS^%?G7MmFF'm7qD)M=J#S.2G&"

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password + SECRET))
	return fmt.Sprintf("%x", h.Sum(nil))
}
