package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

func HashHmac(data string) string {
	h := hmac.New(sha256.New, []byte(os.Getenv("APP_KEY")))

	h.Write([]byte(data))

	/*
	 * The RawURLEncoding encoding uses the URL and filename-safe base64 alphabet,
	 * which replaces the + and / characters with - and _ respectively,
	 * and omits the padding = character at the end of the string.
	 */
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
