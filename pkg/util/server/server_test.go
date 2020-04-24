package server_test

import (
	"go-blog/pkg/util/server"
	"testing"
)

func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
