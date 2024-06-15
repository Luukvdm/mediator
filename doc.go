// Package mediator implements the mediator behavioral pattern.
// This package is heavily inspired by the CSharp implementation https://github.com/jbogard/MediatR build by JBogard.
//
// [Mediator] can be used to [Send] a [Request].
// The [Request] passes through a [Pipeline] before eventually being handled.
//
// The [Mediator] object can also be used to [Publish] and [Subscribe] to notifications.
// It does this without using the reflect package.
package mediator
