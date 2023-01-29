package datacrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type crypt struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	label      string
}

func New(options ...func(*crypt) error) (*crypt, error) {
	c := &crypt{}
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, errors.Wrap(err, "failed to create crypt")
		}
	}
	return c, nil
}

func WithPublicKey(path string) func(*crypt) error {
	return func(c *crypt) error {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to read file with public key - %s", path))
		}

		block, _ := pem.Decode(data)
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse public key - %+v", block))
		}

		c.publicKey = key

		return nil
	}
}

func WithLabel(l string) func(*crypt) error {
	return func(c *crypt) error {
		c.label = l

		return nil
	}
}

func WithPrivateKey(path string) func(*crypt) error {
	return func(c *crypt) error {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to read file with private key - %s", path))
		}

		block, _ := pem.Decode(data)
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse private key - %+v", block))
		}

		c.privateKey = key

		return nil
	}
}

func (c crypt) GetDecryptMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			decryptedBody, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, body, []byte(c.label))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = ioutil.NopCloser(bytes.NewReader(decryptedBody))
			r.ContentLength = int64(len(decryptedBody))

			next.ServeHTTP(w, r)
		})
	}
}
