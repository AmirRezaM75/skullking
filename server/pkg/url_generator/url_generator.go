package url_generator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

type UrlGenerator struct {
	SecretKey string
	BaseURL   string
}

func NewUrlGenerator() UrlGenerator {
	return UrlGenerator{
		// TODO: Not good practice to access env in packages
		SecretKey: os.Getenv("APP_KEY"),
		BaseURL:   os.Getenv("APP_URL"),
	}
}

func (ug UrlGenerator) TemporarySignedRoute(path string, parameters map[string]string, expiration time.Time) (string, error) {
	err := ug.ensureSignedRouteParametersAreNotReserved(parameters)

	if err != nil {
		return "", err
	}

	parameters["expires"] = strconv.FormatInt(expiration.Unix(), 10)

	u, err := ug.buildURL(path, parameters)

	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(ug.SecretKey))

	h.Write([]byte(u.String()))

	parameters["signature"] = hex.EncodeToString(h.Sum(nil))

	u, err = ug.buildURL(path, parameters)

	if err != nil {
		return "", err
	}

	return u.String(), nil

}

func (ug UrlGenerator) ensureSignedRouteParametersAreNotReserved(parameters map[string]string) error {
	message := "'%s' is a reserved parameter when generating signed routes. Please rename your route parameter"

	if _, exists := parameters["signature"]; exists {
		return errors.New(fmt.Sprintf(message, "signature"))
	}

	if _, exists := parameters["expires"]; exists {
		return errors.New(fmt.Sprintf(message, "expires"))
	}

	return nil
}

func (ug UrlGenerator) buildURL(path string, parameters map[string]string) (*url.URL, error) {

	u, err := ug.joinPath(path)

	if err != nil {
		return nil, err
	}

	q := ug.buildQuery(parameters)

	u.RawQuery = q.Encode()

	return u, nil

}

func (ug UrlGenerator) buildQuery(parameters map[string]string) url.Values {
	keys := make([]string, 0, len(parameters))

	for key := range parameters {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	q := url.Values{}

	for _, k := range keys {
		q.Add(k, parameters[k])
	}

	return q
}

func (ug UrlGenerator) joinPath(p string) (*url.URL, error) {
	u, err := url.Parse(ug.BaseURL)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, p)

	return u, nil
}

func (ug UrlGenerator) HasValidSignature(u *url.URL) bool {
	return ug.hasCorrectSignature(u) && ug.signatureHasNotExpired(u)
}

func (ug UrlGenerator) hasCorrectSignature(u *url.URL) bool {
	q := u.Query()
	signature := q.Get("signature")
	q.Del("signature")
	u.RawQuery = q.Encode()

	h := hmac.New(sha256.New, []byte(ug.SecretKey))

	// TODO: Concatenation is not reasonable here
	h.Write([]byte(os.Getenv("APP_URL") + u.String()))

	return signature == hex.EncodeToString(h.Sum(nil))

}

func (ug UrlGenerator) signatureHasNotExpired(u *url.URL) bool {
	expires := u.Query().Get("expires")

	timestamp, _ := strconv.Atoi(expires)

	return expires != "" && int64(timestamp) > time.Now().Unix()
}