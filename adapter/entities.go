package adapter

import (
	"bytes"
)

type Client interface {
	CreateServices(*bytes.Buffer) []Service
	GetService(string) Service
	DeleteService(string) error
	FetchMetadata() (interface{}, error)
}

type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}
