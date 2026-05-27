package api

import "net/http"

type Service interface {
	// Test Handler
	Test() http.HandlerFunc
}

type service struct{}

func NewService() Service {
	return &service{}
}
