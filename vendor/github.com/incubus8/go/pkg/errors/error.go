// Package errors provides the common infrastructure for managing errors.  It
// is primarily a wrapper around github.com/pkg/errors package.  See
// https://godoc.org/github.com/pkg/errors for usage instructions.
//
// While it is possible for most packages to use this package without needing
// the underlying package it is wrapping, certain usecases (such as interacting
// with the recorded stack trace) cannot avoid leakage.

package errors

import (
	"github.com/pkg/errors"
)

type ErrorReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIError is the error type usually returned by functions in the server
// It describes the code, message, and recoded time when an error occurred.
type APIError struct {
	ErrorReason ErrorReason `json:"error"`
	StatusCode  int         `json:"status_code"`
	RecordedAt  string      `json:"recorded_at"`
	Err         error       `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}

	er := e.ErrorReason
	s := ""
	if er.Code != "" {
		s = "API " + er.Code
	}

	if er.Message != "" {
		s = s + ":" + er.Message
	}

	return s
}

type ValidationErrorReason struct {
	Code    string   `json:"code"`
	Message []string `json:"message"`
}

type ValidationError struct {
	ValidationErrorReason ValidationErrorReason `json:"error"`
	StatusCode            int                   `json:"status_code"`
	RecordedAt            string                `json:"recorded_at"`
}

// StackTracer represents a type (usually an error) that can provide a stack
// trace.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Cause returns the underlying cause of the error, if possible.  See
// https://godoc.org/github.com/pkg/errors#Cause for further details.
func Cause(err error) error {
	return errors.Cause(err)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error. See
// https://godoc.org/github.com/pkg/errors#Errorf for further details
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// New returns an error with the supplied message. See
// https://godoc.org/github.com/pkg/errors#New for further details
func New(message string) error {
	return errors.New(message)
}

// Wrap returns an error annotating err with message. If err is nil, Wrap
// returns nil.  See https://godoc.org/github.com/pkg/errors#Wrap for more
// details.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with the format specifier. If err is
// nil, Wrapf returns nil. See https://godoc.org/github.com/pkg/errors#Wrapf
// for more details.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
