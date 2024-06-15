//go:build generate
// +build generate

package mediator

import (
	_ "github.com/golangci/golangci-lint"
	_ "github.com/vektra/mockery/v2"
)

//go:generate rm -rf mocks
//go:generate go run github.com/vektra/mockery/v2
