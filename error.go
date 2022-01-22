// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"bytes"
	"encoding/json"
	"io"
)

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
