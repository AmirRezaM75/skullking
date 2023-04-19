package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"
)

type UrlGenerator struct {
	SecretKey string
	BaseURL   string
}

func (ug UrlGenerator) TemporarySignedRoute(path string, parameters map[string]string, expiration time.Time) (string, error) {
	err := ug.ensureSignedRouteParametersAreNotReserved(parameters)

	if err != nil {
		return "", err
	}

	parameters["expires"] = strconv.FormatInt(expiration.Unix(), 10)

	p := ksort(parameters)

	u, err := ug.url(path, p)

	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(ug.SecretKey))

	h.Write([]byte(u.String()))

	p["signature"] = hex.EncodeToString(h.Sum(nil))

	u, err = ug.url(path, p)

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

func (ug UrlGenerator) url(path string, parameters map[string]string) (*url.URL, error) {

	u, err := ug.joinPath(path)

	if err != nil {
		return nil, err
	}

	q := ug.buildQuery(parameters)

	u.RawQuery = q.Encode()

	return u, nil

}

func (ug UrlGenerator) buildQuery(parameters map[string]string) url.Values {
	q := url.Values{}

	for k, v := range parameters {
		q.Add(k, v)
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
