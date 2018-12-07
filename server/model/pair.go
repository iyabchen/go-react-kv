package model

import (
	"fmt"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Pair represents a key-value pair.
type Pair struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewPair validate the input and create the pair.
func NewPair(key string, value string) (*Pair, error) {
	k := strings.TrimSpace(key)
	v := strings.TrimSpace(value)
	r, err := regexp.Compile("[a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	if !r.MatchString(k) {
		return nil, fmt.Errorf("Key %s contains non alpha-numeric characters",
			k)
	}
	if !r.MatchString(v) {
		return nil, fmt.Errorf("Value %s contains non alpha-numeric characters",
			v)
	}

	return &Pair{
		ID:    uuid.NewV4().String(),
		Key:   k,
		Value: v,
	}, nil
}
