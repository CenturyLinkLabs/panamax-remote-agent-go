package adapter

import (
	"bytes"
)

//TODO: rename
type AdapterClient interface {
	CreateServices(*bytes.Buffer) []Service
	GetService(string) Service
	DeleteService(string) error
	FetchMetadata() (interface{}, error)
}

type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}
