package types

import "strings"

type QueryResExample struct {
	Value string `json:"value"`
}

func (r QueryResExample) String() string {
	return r.Value
}

type QueryResNames []string

func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}
