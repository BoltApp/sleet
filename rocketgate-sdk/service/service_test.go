package service

import (
	"fmt"
	"testing"
)

func TestNewGatewayResponse(t *testing.T) {
	service := NewGatewayService()
	fmt.Println(service)
}
