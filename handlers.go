// Copyright 2015 The ATC Authors. All Rights Reserved.
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
// more infomation: http://ath.wiki
package atc

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"encoding/json"
	"github.com/adolphlxm/atc/context"
	"github.com/adolphlxm/atc/rpc/thrift/gen/atcrpc"
	"net/url"
	"path"
	"strings"
)

var (
	HTTPMETHODS = map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"DELETE":  "DELETE",
		"PATCH":   "PATCH",
		"OPTIONS": "OPTIONS",
		"HEAD":    "HEAD",
		"TRACE":   "TRACE",
		"CONNECT": "CONNECT",
	}
)

// HandlerInterface is an interface to uniform all controller handler.
type HandlerInterface interface {
	Init(ctx *context.Context)
	Get()
	Post()
	Delete()
	Put()
	Patch()
	Head()
	Options()
	Websocket()
}

type Handler struct {
	// The current context interface
	Ctx *context.Context

	// The current Conn represents a WebSocket connection.
	//ws *websocket.Conn
}

// Init generates default values of controller operations.
func (h *Handler) Init(ctx *context.Context) {
	h.Ctx = ctx
}

// Get
func (h *Handler) Get() {
	h.Error404(404)
}

// Post
func (h *Handler) Post() {
	h.Error404(404)
}

// Delete
func (h *Handler) Delete() {
	h.Error404(404)
}

// Put
func (h *Handler) Put() {
	h.Error404(404)
}

// PATCH
func (h *Handler) Patch() {
	h.Error404(404)
}

// Head
func (h *Handler) Head() {
	h.Error404(404)
}

// Options
func (h *Handler) Options() {
	h.Error404(404)
}

// Websocket
func (h *Handler) Websocket() {
	for {
		input, err := h.Recevie()
		if err != nil {
			h.Error406(407)
			continue
		}

		h.Ctx.Reset()
		h.Ctx.ReqType = context.RPC_WEBSOCKET

		//base64
		enbyte, err := base64.StdEncoding.DecodeString(input.Body)
		if err != nil {
			h.Error406(408)
			continue
		}

		body, err := url.ParseQuery(string(enbyte))

		if err != nil {
			h.Error406(409)
			continue
		}
		h.Ctx.Request.Form = body

		if HTTPMETHODS[input.Method] == "" {
			h.Error406(410)
		}
		if input.Module == "" {
			h.Error406(411)
		}

		ModuleName := strings.Split(input.Module, ".")
		ExecuteHandler(input.Method, path.Join(input.Version, path.Join(ModuleName...)), h)
	}
}

// TODO Thrift
func (h *Handler) CallBack(a *atcrpc.ReqHandler, body map[string]string) (r string, err error) {
	var content []byte

	moduleName := strings.Split(a.Handler, ".")
	h.Ctx = context.NewContext(nil, nil)
	h.Ctx.ReqType = context.RPC_THRIFT

	ExecuteHandler(a.Method, path.Join(a.Version, path.Join(moduleName...)), h)
	content, err = json.Marshal(h.Ctx.Data())
	r = string(content)
	return
}

type JsonRPC struct {
	// Request API version
	// e.g.
	//	V1,V2...
	Version string `json:"version"`
	// Request method
	// e.g.
	//	GET,POST,PUT...
	Method string `json:"method"`
	// Request module
	// e.g.
	//	user.login
	Module string `json:"module"`
	// Request body Encoding agnostic text or binary string
	Body string `json:"body"`
}

// Receive receives single frame from ws, unmarshaled by cd.Unmarshal and stores in v.
func (h *Handler) Recevie() (*JsonRPC, error) {
	var reply *JsonRPC
	err := websocket.JSON.Receive(h.Ctx.WS, &reply)
	h.Ctx.Reset()
	return reply, err
}

// Send sends v marshaled by cd.Marshal as single frame to ws.
func (h *Handler) Send() error {
	err := websocket.JSON.Send(h.Ctx.WS, h.Ctx.Data())
	return err
}

// JSON writes json to response body.
func (h *Handler) JSON() error {
	err := h.Ctx.SaveJSON(h.Ctx.Data())
	if err != nil {
		h.Error406(406)
	}
	return err
}

// Use:
//	1. this.Error400(400).OutPut()
//	2. this.Error400(400).Message("Error messages").OutPut()
func (h *Handler) Error400(code int) *Error {
	return h.Error(400, code)
}

func (h *Handler) Error401(code int) *Error {
	return h.Error(401, code)
}

func (h *Handler) Error403(code int) *Error {
	return h.Error(403, code)
}

func (h *Handler) Error404(code int) *Error {
	return h.Error(404, code)
}

func (h *Handler) Error406(code int) *Error {
	return h.Error(406, code)
}

func (h *Handler) Error500(code int) *Error {
	return h.Error(500, code)
}

func (h *Handler) Error(status, code int) *Error {
	response := NewError(h.Ctx)
	return response.Code(status, code)
}
