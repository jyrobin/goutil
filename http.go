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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	neturl "net/url"
	"strings"
)

var httpVerbs = StrMap{}

func init() {
	httpVerbs[http.MethodGet] = http.MethodGet
	httpVerbs[http.MethodPost] = http.MethodPost
	httpVerbs[http.MethodPut] = http.MethodPut
	httpVerbs[http.MethodPatch] = http.MethodPatch
	httpVerbs[http.MethodDelete] = http.MethodDelete
}

const (
	ctAppJson = "application/json"
)

func HandleError(w http.ResponseWriter, err error, statusCode int) bool {
	if err == nil {
		return false
	}
	http.Error(w, string(err.Error()), statusCode)
	return true
}

func ReadResponseBody(res *http.Response) ([]byte, error) {
	if res == nil || res.Body == nil {
		return nil, fmt.Errorf("Nil response")
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func ReadRequestBody(req *http.Request) ([]byte, error) {
	if req == nil || req.Body == nil {
		return nil, fmt.Errorf("Nil response")
	}

	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}

// Marshal/unmarshal all mean json, all ignore whether Content-Type is applicattion/json
// - maybe check if content-type is application/json later
func UnmarshalResponse(res *http.Response, ret interface{}) error {
	buf, err := ReadResponseBody(res)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("StatusCode %d", res.StatusCode)
	}
	if len(buf) == 0 {
		return fmt.Errorf("Empty body")
	}
	return json.Unmarshal(buf, ret)
}

func UnmarshalRequest(req *http.Request, ret interface{}) error {
	buf, err := ReadRequestBody(req)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return fmt.Errorf("Empty body")
	}
	return json.Unmarshal(buf, ret)
}

// REF: https://gist.github.com/rjz/fe283b02cbaa50c5991e1ba921adf7c9

// Failure should yield an HTTP 415 (`http.StatusUnsupportedMediaType`)
func HasContentType(r *http.Request, mimetype string) bool {
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

// REF: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func DecodeRequestBody(req *http.Request, strict bool, ret interface{}) error {
	if strict && !HasContentType(req, ctAppJson) {
		msg := "Content-Type header is not application/json"
		return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
	}

	var body io.Reader = req.Body
	if strict {
		body = io.LimitReader(req.Body, 1048576) // 1MB
	}

	dec := json.NewDecoder(body)
	//if strict {
	//	dec.DisallowUnknownFields()
	//}

	defer req.Body.Close()

	err := dec.Decode(&ret)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

// requests

func HttpRequest(method, url, contentType string, data []byte, headers StrMap) (*http.Request, error) {
	mthd, ok := httpVerbs[strings.ToUpper(method)]
	if !ok {
		return nil, fmt.Errorf("%s not supported", method)
	}

	var reader io.Reader = nil
	if data != nil && mthd != http.MethodGet {
		reader = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(mthd, url, reader)
	if err != nil {
		return nil, err
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func GetRequest(url, contentType string, headers StrMap) (*http.Request, error) {
	return HttpRequest(http.MethodGet, url, contentType, nil, headers)
}

func AjaxRequest(method, url string, data []byte, headers StrMap) (*http.Request, error) {
	return HttpRequest(method, url, ctAppJson, data, headers)
}

func AjaxGetRequest(url string, headers StrMap) (*http.Request, error) {
	return AjaxRequest(http.MethodGet, url, nil, headers)
}

func AjaxPostRequest(url string, jsonStr []byte, headers StrMap) (*http.Request, error) {
	return AjaxRequest(http.MethodPost, url, jsonStr, headers)
}

func AjaxPutRequest(url string, jsonStr []byte, headers StrMap) (*http.Request, error) {
	return AjaxRequest(http.MethodPut, url, jsonStr, headers)
}

func FormRequest(method, url string, data, headers StrMap) (*http.Request, error) {
	mthd, ok := httpVerbs[strings.ToUpper(method)]
	if !ok || mthd == http.MethodGet {
		return nil, fmt.Errorf("%s not supported", method)
	}

	form := neturl.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	req, err := http.NewRequest(mthd, url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	for h, v := range headers {
		req.Header.Add(h, v)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// http calls

func HttpGet(url string) ([]byte, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func HttpDo(method, url, contentType string, data []byte, headers StrMap) (*http.Response, error) {
	req, err := HttpRequest(method, url, contentType, data, headers)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func HttpCall(method, url, contentType string, data []byte, headers StrMap) ([]byte, *http.Response, error) {
	res, err := HttpDo(method, url, contentType, data, headers)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return buf, res, nil
}

func Http(method, url, contentType string, data []byte, headers StrMap) ([]byte, error) {
	buf, _, err := HttpCall(method, url, contentType, data, headers)
	return buf, err
}

func HttpSend(method, url string, data []byte, headers StrMap) ([]byte, error) {
	return Http(method, url, "", data, headers)
}

func Ajax(method, url string, jsonStr []byte, headers StrMap) ([]byte, error) {
	return Http(method, url, ctAppJson, jsonStr, headers)
}

func AjaxGet(url string, headers StrMap) ([]byte, error) {
	return Http(http.MethodGet, url, ctAppJson, nil, headers)
}

func AjaxUnmarshal(method, url string, jsonStr []byte, headers StrMap, ret interface{}) error {
	res, err := HttpDo(method, url, ctAppJson, nil, headers)
	if err != nil {
		return err
	}
	return UnmarshalResponse(res, ret)
}

func AjaxGetUnmarshal(url string, headers StrMap, v interface{}) error {
	return AjaxUnmarshal(http.MethodGet, url, nil, headers, v)
}

func AjaxPost(url string, jsonStr []byte, headers StrMap) ([]byte, error) {
	return Ajax(http.MethodPost, url, jsonStr, headers)
}

func AjaxPostUnmarshal(url string, jsonStr []byte, headers StrMap, v interface{}) error {
	return AjaxUnmarshal(http.MethodPost, url, jsonStr, headers, v)
}

func AjaxPut(url string, jsonStr []byte, headers StrMap) ([]byte, error) {
	return Ajax(http.MethodPut, url, jsonStr, headers)
}

func AjaxPutUnmarshal(url string, jsonStr []byte, headers StrMap, v interface{}) error {
	return AjaxUnmarshal(http.MethodPut, url, jsonStr, headers, v)
}

// simple ajax calls, sending and receiving Map

func SimpleAjaxRequest(method, url string, data Map, headers StrMap) (*http.Request, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return AjaxRequest(method, url, jsonStr, headers)
}

func SimpleAjaxDo(method, url string, data Map, headers StrMap) (*http.Request, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return AjaxRequest(method, url, jsonStr, headers)
}

func SimpleAjax(method, url string, data Map, headers StrMap) (Map, error) {
	var jsonStr []byte = nil
	var err error
	if data != nil && method != http.MethodGet {
		if jsonStr, err = json.Marshal(data); err != nil {
			return nil, err
		}
	}

	var ret Map
	if err = AjaxUnmarshal(method, url, jsonStr, headers, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SimpleAjaxGet(url string, headers StrMap) (Map, error) {
	return SimpleAjax(http.MethodGet, url, nil, headers)
}

func SimpleAjaxPost(url string, data Map, headers StrMap) (Map, error) {
	return SimpleAjax(http.MethodPost, url, data, headers)
}

func SimpleAjaxPut(url string, data Map, headers StrMap) (Map, error) {
	return SimpleAjax(http.MethodPut, url, data, headers)
}

func SimpleAjaxUnmarshal(method, url string, data Map, headers StrMap, ret interface{}) error {
	var jsonStr []byte = nil
	var err error
	if data != nil && method != http.MethodGet {
		if jsonStr, err = json.Marshal(data); err != nil {
			return err
		}
	}
	return AjaxUnmarshal(method, url, jsonStr, headers, ret)
}

// for calling ServeHTTP directly

type responseWriter struct {
	header     http.Header
	body       *bytes.Buffer
	statusCode int
}

func NewResponseWriter() *responseWriter {
	return &responseWriter{
		header:     make(http.Header),
		body:       new(bytes.Buffer),
		statusCode: 200,
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	w.body.Write(buf)
	return len(buf), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func (w *responseWriter) Bytes() []byte {
	return w.body.Bytes()
}

func (w *responseWriter) String() string {
	return w.body.String()
}

func (w *responseWriter) Unmarshal(ret interface{}) error {
	if w.statusCode >= 400 {
		return fmt.Errorf("StatusCode %d", w.statusCode)
	}
	return json.Unmarshal(w.Bytes(), ret)
}
