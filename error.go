// Copyright (c) 2021 Jing-Ying Chen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
