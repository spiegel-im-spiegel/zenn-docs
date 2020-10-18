package encode

import (
	"crypto"
	"errors"
	"io"
)

var (
	ErrNoImplement = errors.New("no implementation")
)

//Value returns hash value string from io.Reader
func Value(r io.Reader, alg crypto.Hash) ([]byte, error) {
	if !alg.Available() {
		return nil, ErrNoImplement
	}
	h := alg.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
