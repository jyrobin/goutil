// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	InternalServerError = 500
)

type Error struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []ErrorItem `json:"errors,omitempty"`
}

type ErrorItem struct {
	Message string `json:"message,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

func StdError(code int, msg string, items ...ErrorItem) Error {
	return Error{code, msg, items}
}
func NotFoundError(msgs ...string) Error {
	msg := "Not Found"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return StdError(NotFound, msg)
}
func IsError(err error, code int) bool {
	e, ok := err.(Error)
	return ok && e.Code == code
}
func ErrorCode(err error, otherwise int) int {
	if e, ok := err.(Error); ok {
		return e.Code
	}
	return otherwise
}

func IsNotFound(err error) bool {
	return IsError(err, NotFound)
}
func IsUnauthorized(err error) bool {
	return IsError(err, NotFound)
}

// BufferError is not really of an error type but to "tunnel" return data
// from functions that only return errors
type BufferError interface {
	error
	io.ReadWriter
	Len() int
	Bytes() []byte
	String() string
}

type BytesBufferError struct {
	bytes.Buffer
	msg string
}

func NewBufferError(msg string, data []byte) *BytesBufferError {
	return &BytesBufferError{
		Buffer: *bytes.NewBuffer(append([]byte(nil), data...)),
		msg:    msg,
	}
}

func NewStringError(msg, data string) *BytesBufferError {
	return &BytesBufferError{
		Buffer: *bytes.NewBufferString(data),
		msg:    msg,
	}
}

func NewJsonError(msg string, data interface{}) *BytesBufferError {
	buf, err := json.Marshal(data)
	if err != nil {
		return NewBufferError(err.Error(), nil)
	}
	return NewBufferError(msg, buf)
}

func (e BytesBufferError) Error() string { return e.msg }
