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
// more infomation: http://atc.wiki
package atc

import (
	"bytes"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/adolphlxm/atc/context"
	"github.com/adolphlxm/atc/logs"
	"fmt"
)

type Location int

const (
	BEFORE_ROUTE   Location = iota //匹配路由之前
	BEFORE_HANDLER                 //匹配到路由后,执行Handler之前
	AFTER                          //执行完所有逻辑后
)

// FilterFunc defines a filter function which is invoked before the controller handler is executed.
type FilterFunc func(*context.Context)

// A HandlerRouter store routing rules,
// routers and filters.
type HandlerRouter struct {
	mu sync.RWMutex
	// A list of routes
	routers []*Router
	// A list of handler filters,support BEFORE_ROUTE,BEFORE_HANDLER,AFTER
	filters map[Location][]*Router
}

// NewHandlerRouter returns a new HandlerRouter.
func NewHandlerRouter() (*HandlerRouter, error) {
	a := &HandlerRouter{
		filters: make(map[Location][]*Router),
	}
	//http.Handle("/", a)
	return a, nil
}

// AddRouter returns a point to the Router
//
// RESTful usage:
// 	AddRouter("/V1/user/{userid:[0-9]?}",&user.UserHandler{})
// 	AddRouter("/V1/user/group",&user.GroupHandler{})
func (h *HandlerRouter) AddRouter(pattern string, c HandlerInterface) *Router {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create a new route
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()

	// Check the routing legal
	for _, r := range h.routers {
		if pattern == r.Pattern {
			panic("ATC AddRouter: new router failed for pattern " + pattern + ": routing repeated!")
		}
	}

	// Register route
	router, err := newRouter(pattern, t)

	// TODO
	// Register thrift
	//if Aconfig.ThriftSupport {
	//	processor := atcrpc.NewAtcrpcThriftProcessor(c)
	//	Thrift.RegisterProcessor(pattern, processor)
	//}

	if err != nil {
		panic("ATC AddRouter: new router failed for pattern " + pattern + ":" + err.Error())
	}

	h.routers = append(h.routers, router)
	return router
}

// AddFilter is to add filters
func (h *HandlerRouter) AddFilter(location Location, pattern string, f FilterFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create a new route
	router, err := newRouter(pattern, nil)
	if err != nil {
		panic("ATC AddFilter: new router failed for filter " + pattern + ":" + err.Error())
	}
	router.RunFilter = f

	h.filters[location] = append(h.filters[location], router)
}

// Handler is a simple interface to a http.Handler browser client.
// It checks if Origin header is valid URL by default.
func (h *HandlerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := path.Clean(r.URL.Path)
	//Logger.Warn(requestPath)
	if !strings.HasPrefix(requestPath, "/") {
		requestPath = "/" + requestPath
		r.URL.Path = requestPath
	}

	ctx := context.NewContext(w, r)

	defer h.recoverPanic(ctx)

	logs.Trace("%s %s for %v", r.Method, r.URL.Path, ctx.IP())

	// Static files routing.
	if Aconfig.FrontSupport {
		if err := frontStaticRouter(ctx); err == nil {
			return
		}
	}

	if ctx.IsWebsocket() {
		// ServeHTTP implements the http.Handler interface for a WebSocket
		websocket.Handler(func(ws *websocket.Conn) {
			// Override default Read/Write timeout with sane value for a web socket request
			ws.SetDeadline(time.Now().Add(time.Hour * 24))

			ctx.WS = ws
			h.findRoute("WS", requestPath, ctx)
		}).ServeHTTP(w, r)

	} else {
		// RESTFUL handler
		h.findRoute(ctx.Method(), requestPath, ctx)
		ctx.MultipartFormMaxMemory(Aconfig.PostMaxMemory)
	}

	// Exit handler
	if ctx.GetStatus() != http.StatusOK {
		return
	}
	h.findFilter(AFTER, requestPath, ctx)
}

func (h *HandlerRouter) ExecuteHandler(method, requestPath string, c *Handler) {
	defer h.recoverPanic(c.Ctx)

	h.findRoute(method, requestPath, c.Ctx)
	// Exit handler
	if c.Ctx.GetStatus() != http.StatusOK {
		return
	}
	h.findFilter(AFTER, requestPath, c.Ctx)
}

func (h *HandlerRouter) recoverPanic(c *context.Context) {
	if err := recover(); err != nil {
		// Is open panic
		if !Aconfig.Debug {
			logs.Fatal("%s request recover: %v", c.Path(), err)
		}

		logs.Error("%s request recover: %v", c.Path(), err)
	}
}

// finds the matching route given a cleaned path
func (h *HandlerRouter) findRoute(method, requestPath string, c *context.Context) {
	error := NewError(c)

	r := h.matchRouter(requestPath)
	if r == nil {
		error.Code(404, 404).JSON()
		return
	}
	c.RunHandler = r.HandlerType

	h.findFilter(BEFORE_ROUTE, requestPath, c)
	// Exit handler
	if c.GetStatus() != http.StatusOK {
		return
	}

	if r != nil {
		// Loading controller handler before the filter
		// If the HTTP status code is not 200, stop running,
		// apply to websocket.
		h.findFilter(BEFORE_HANDLER, requestPath, c)
		// Exit handler
		if c.GetStatus() != http.StatusOK {
			return
		}
		switch c.ReqType {
		case context.RPC_HTTP:
			c.SetParams(r.MatchParams(requestPath))
		}

		vc := reflect.New(r.HandlerType)
		execController, ok := vc.Interface().(HandlerInterface)
		if !ok {
			error.Code(500, 500).JSON()
		}

		execController.Init(c)
		switch method {
		case "GET":
			execController.Get()
		case "POST":
			execController.Post()
		case "DELETE":
			execController.Delete()
		case "PUT":
			execController.Put()
		case "PATCH":
			execController.Patch()
		case "HEAD":
			execController.Head()
		case "OPTIONS":
			execController.Options()
		case "WS":
			execController.Websocket()
		default:
			execController.Get()
		}
	}
}

func (h *HandlerRouter) matchRouter(requestPath string) *Router {
	if requestPath == "" {
		return nil
	}

	for _, r := range h.routers {
		if r.MatchPath(requestPath) {
			return r
		}
	}
	return nil
}

// finds the matching filter given a cleaned path
func (h *HandlerRouter) findFilter(location Location, requestPath string, c *context.Context) {
	if r, ok := h.filters[location]; ok {
		for _, filter := range r {
			// check method and path
			if filter.MatchPath(requestPath) {
				filter.RunFilter(c)
				logs.Trace("Execution handler filter path:%v", filter.Pattern)
			}
		}
	}
}

type Router struct {
	// A simple HTTP handler ATC
	RunFilter FilterFunc

	// Routing patterns
	Pattern string

	// Match the routing the regular of success
	Regexp *regexp.Regexp

	// Match the keys from the pattern
	Params []string

	HandlerType reflect.Type
}

func newRouter(pattern string, t reflect.Type) (r *Router, err error) {
	r = &Router{
		HandlerType: t,
		Pattern:     pattern,
	}

	//Check regexp
	err = r.regexpRouter()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Router) MatchPath(path string) bool {
	if r.Pattern == path {
		return true
	} else if r.Regexp != nil {
		if r.Regexp.MatchString(path) {
			return true
		}
	}
	return false
}

func (r *Router) MatchParams(path string) map[string]string {
	params := make(map[string]string)
	if r.Regexp == nil {
		return params
	}

	if regstr := r.Regexp.FindStringSubmatch(path); len(regstr) > 0 {
		for i, match := range regstr[1:] {
			params[r.Params[i]] = match
		}
	}
	return params
}

func (r *Router) regexpRouter() (err error) {
	metaPattern := regexp.QuoteMeta(r.Pattern)
	if metaPattern != r.Pattern {

		//Create a buffer
		exprPattern := bytes.NewBufferString("^")
		defaultPattern := "[^/]+"
		idxs, err := braceIndices(r.Pattern)
		if err != nil {
			return err
		}

		var end int
		for i := 0; i < len(idxs); i += 2 {
			raw := regexp.QuoteMeta(r.Pattern[end:idxs[i]])
			end = idxs[i+1]
			parts := strings.SplitN(r.Pattern[idxs[i]+1:end-1], ":", 2)
			patt := defaultPattern
			if len(parts) == 2 {
				patt = parts[1]
			}
			if patt == "" {
				return fmt.Errorf("router: missing name or pattern in %q",
					r.Pattern[idxs[i]:end])
			}

			switch patt[len(patt)-1:] {
			case "?","*":
				switch path.Clean(raw) {
				case "/":
					raw = "([/]?)"
				default:
					raw = path.Clean(raw) + "([/]?)"
				}
				r.Params = append(r.Params, "")
			}

			// Build the regexp parameters.
			r.Params = append(r.Params, parts[0])

			// Build the regexp pattern.
			fmt.Fprintf(exprPattern, "%s(%s)", raw, patt)
		}

		exprPattern.WriteByte('$')

		//Compile parses a regular expression and returns, if successful
		r.Regexp, err = regexp.Compile(exprPattern.String())
	}
	return err
}

func braceIndices(s string) ([]int, error) {
	var level, idx int
	var idxs []int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("router: unbalanced braces in %q", s)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("router: unbalanced braces in %q", s)
	}
	return idxs, nil
}