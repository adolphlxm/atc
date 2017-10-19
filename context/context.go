// Copyright 2016 ATC Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ATC is an open-source, automated test framework for the Go programming language.
//
// more infomation: http://atc.wiki
package context

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"code.google.com/p/go.net/websocket"
	"reflect"
)

var (
	acceptsJSONRegex = regexp.MustCompile(`(application/vnd.atc+json)(?:,|$)`)
)

type RequestType byte

const (
	RPC_HTTP = RequestType(iota)
	RPC_WEBSOCKET
	RPC_THRIFT
)

type Context struct {
	// The current response writer
	ResponseWriter http.ResponseWriter

	// The current request, data and body.
	Request *http.Request

	// The current Conn represents a WebSocket connection.
	WS *websocket.Conn

	RunHandler reflect.Type

	ReqType RequestType

	// The HTTP response header with status code.
	// If it is 200 as normal,
	// else the error code.
	status int

	// The current route data stored in a map
	params map[string]string

	// Arbitrary user data stored in a map
	data map[string]interface{}

	// ParseMultipartForm parses a request body as multipart/form-data.
	// The whole request body is parsed and up to a total of maxMemory bytes
	maxMemory int64
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		ReqType:        RPC_HTTP,
		status:         http.StatusOK,
		params:         make(map[string]string),
		data:           make(map[string]interface{}),
	}
}

func (ctx *Context) Reset() {
	ctx.ReqType = RPC_HTTP
	ctx.status = http.StatusOK
	ctx.params = nil
	ctx.data = nil
}

/************************************/
/********** Request input ***********/
/************************************/

// Header returns request header item string by a given string
func (ctx *Context) Header(key string) string {
	return ctx.Request.Header.Get(key)
}

// Path returns the path for the request
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

// Scheme returns request scheme as "http" or "https".
func (ctx *Context) Scheme() string {
	if ctx.Request.URL.Scheme != "" {
		return ctx.Request.URL.Scheme
	}
	if ctx.Request.TLS == nil {
		return "http"
	}
	return "https"
}

// Method returns http request method.
func (ctx *Context) Method() string {
	return ctx.Request.Method
}

// IsMethod returns request method is boolean
// usage:
// IsMethod("GET")
// IsMethod("POST")
// IsMethod("PUT")
// IsMethod("DELETE")
func (ctx *Context) IsMethod(method string) bool {
	return ctx.Method() == method
}

// IsWebsocket returns boolean of this request is in webSocket.
func (ctx *Context) IsWebsocket() bool {
	return ctx.Header("Upgrade") == "websocket"
}

// AcceptsJSON Checks if request accepts json response
func (ctx *Context) AcceptsJSON() bool {
	return acceptsJSONRegex.MatchString(ctx.Header("Accept"))
}

// IP returns request ip
func (ctx *Context) IP() string {
	address := ctx.Header("X-Real-IP")
	if len(address) > 0 {
		return address
	}

	address = ctx.Header("X-Forwarded-For")
	if len(address) > 0 {
		return address
	}

	return ctx.Request.RemoteAddr
}

// Query returns input data item string by a given string.
func (ctx *Context) Query(key string) string {
	if val := ctx.Param(key); val != "" {
		return val
	}
	if err := ctx.parseForm(); err != nil {
		return ""
	}
	return ctx.Request.Form.Get(key)
}

func (ctx *Context) Param(key string) string {
	if ctx.params == nil {
		ctx.params = make(map[string]string)
	}
	if v, ok := ctx.params[key]; ok {
		return v
	}
	return ""
}

func (ctx *Context) SetParam(key string, val string) {
	if ctx.params == nil {
		ctx.params = make(map[string]string)
	}
	ctx.params[key] = val
}

func (ctx *Context) SetParams(params map[string]string) {
	if ctx.params == nil {
		ctx.params = make(map[string]string)
	}
	ctx.params = params
}

func (ctx *Context) parseForm() (err error) {

	if ctx.Request.Form == nil {
		if strings.Contains(ctx.Header("Context-Type"), "multipart/form-data") {
			// enctype = multipart/form-data
			err = ctx.Request.ParseMultipartForm(ctx.maxMemory)
		} else {
			err = ctx.Request.ParseForm()
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) MultipartFormMaxMemory(maxMemory int64) {
	ctx.maxMemory = maxMemory
}

/************************************/
/********* Response output **********/
/************************************/

// ResHeader sets response header item string via given key.
func (ctx *Context) ResponseHeader(key, val string) {
	if len(val) == 0 {
		ctx.ResponseWriter.Header().Del(key)
	} else {
		ctx.ResponseWriter.Header().Set(key, val)
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (ctx *Context) Write(content []byte) (int, error) {
	// Set header
	ctx.ResponseHeader("Content-Type", "application/vnd.atc+json")

	if ctx.status != 200 {
		ctx.WriteHeader(ctx.status)
	}
	return ctx.ResponseWriter.Write(content)
}

// WriteHeader sends an HTTP response header with status error codes.
func (ctx *Context) WriteHeader(i int) {
	ctx.ResponseWriter.WriteHeader(i)
}

func (ctx *Context) GetStatus() int {
	return ctx.status
}

// SetStatus sets response status error code.
func (ctx *Context) SetStatus(status int) {
	ctx.status = status
}

// returns the stored data for this request.
func (ctx *Context) GetData(k string) interface{} {
	if v, ok := ctx.data[k]; ok {
		return v
	}
	return nil
}

// Data return the implicit data in the input
func (ctx *Context) Data() map[string]interface{} {
	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	return ctx.data
}

// Set saves data for this request
func (ctx *Context) SetData(k string, v interface{}) {
	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	ctx.data[k] = v
}

// Set saves data for this request
func (ctx *Context) SetDatas(data map[string]interface{}) {
	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	ctx.data = data
}

// JSON
func (ctx *Context) SaveJSON(data interface{}) error {
	var (
		err     error
		content []byte
	)

	if ctx.IsWebsocket() {
		err = websocket.JSON.Send(ctx.WS, data)
		//ctx.Reset()
	} else {
		content, err = json.Marshal(data)
		_, err = ctx.Write(content)
	}

	return err
}
