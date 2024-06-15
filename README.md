Mediator
=======

[![GoDoc](https://godoc.org/github.com/luukvdm/mediator?status.svg)](https://pkg.go.dev/github.com/luukvdm/mediator)
[![Coverage](https://img.shields.io/codecov/c/github/luukvdm/mediator)](https://codecov.io/gh/luukvdm/mediator)
[![License](https://img.shields.io/github/license/luukvdm/mediator)](./LICENSE)

A mediator implementation for Golang. 
The goal of this package is to make using CQRS easy and testable in Golang 
using the [mediator design pattern](https://wikipedia.org/wiki/Mediator_pattern).

Mediator supports sending requests and publishing and subscribing to notifications.
Requests and notifications are passed through a pipeline that can be used for things like logging and tracing.
Mediator also makes it easy to mock sending requests. So Http and gRPC handlers can be unit tested.

Examples can be found on [pkg.go.dev](https://pkg.go.dev/github.com/luukvdm/mediator).

Install it with the Go CLI.
```bash
go get github.com/luukvdm/mediator
```