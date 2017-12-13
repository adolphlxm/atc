// Copyright 2015 ATC Author. All Rights Reserved.
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
package atc

import (
	"context"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/adolphlxm/atc/logs"
	"github.com/adolphlxm/atc/grace"
)

// ATC framework version.
const VERSION = "0.8.2"

var Route *RouterGroup
var graceNodeTree *grace.Grace

// Running :
//	1. ORM
//	2. Thrift
//	3. HTTP
func Run() {
	var (
		err error
	)

	// If support Thrift serve.
	if Aconfig.ThriftSupport {
		err = ThriftRPC.Run()
		if err != nil {
			panic(err)
		}
		logs.Trace("Thrift server Running on %v", ThriftRPC.Addr())
	}

	// If support HTTP serve.
	if Aconfig.HTTPSupport {
		HttpAPP.Run()
	}

	logs.Trace("Process PID for %d", os.Getpid())

	stop()
}

// Wait for all HTTP and Thrift fetches to complete.
func stop() {

	// Signal
	//	1. TERM,INT 立即终止
	//	2. QUIT 平滑终止
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	for {
		sig := <-sigChan
		logs.Trace("%v", sig)

		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			os.Exit(1)
		case syscall.SIGQUIT:
			logs.Trace("shutdown: start...")
		}

		// Grace exit.
		if err := graceNodeTree.Stop(); err != nil {
			logs.Error("shutdown: grace exit, err:%s", err.Error())
			continue
		}

		break
	}

	if Aconfig.ThriftSupport {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(Aconfig.ThriftQTimeout)*time.Second)
		ThriftRPC.Shutdown(ctx)
		logs.Trace("shutdown: thrift, biggest waiting for %ds...", Aconfig.ThriftQTimeout)
	}

	if Aconfig.HTTPSupport {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(Aconfig.HTTPQTimeout)*time.Second)
		HttpAPP.Server.Shutdown(ctx)
		logs.Trace("shutdown: http, biggest waiting for %ds...", Aconfig.HTTPQTimeout)
		time.Sleep(1 * time.Millisecond)
	}

	// TODO
}

type RouterGroup struct {
	versionPath string
}

func (group *RouterGroup) Group(versionPath string) *RouterGroup {
	return &RouterGroup{
		versionPath: versionPath,
	}
}

// AddRouter add routing group.
//
// RESTful usage:
//	v1 := atc.Route.Group("V1"){
//		// GET/POST... http://localhost/V1/user/login
//		v1.AddRouter("user",&user.UserHandler{})
//	}
//	v2 := atc.Route.Group("V2"){
//		//	GET/POST... http://localhost/V2/user | GET/POST... http://localhost/V2/user/{id}
//		v2.AddRouter("user.{id:[0-9]?}",&user2.UserHandler{})
//
//		//	GET/POST... http://localhost/V2/user/group
//		v2.AddRouter("user.group",&user2.GroupHandler{})
//	}
//
// Request usage:
// RPC_WEBSOCKET:
//	{"version":"V1","method":"GET(POST...)","handler":"users.login","body":""}
//	{"version":"V2","method":"GET(POST...)","handler":"shop.car.info","body":""}
func (group *RouterGroup) AddRouter(module string, c HandlerInterface) {
	moduleName := strings.Split(module, ".")
	module = path.Join(moduleName...)
	HttpAPP.Handler.AddRouter(path.Join("/", group.versionPath, module), c)
}

// AddFilter add filter group.
//
// RESTful usage:
//	v1 := atc.Route.Group("V1"){
//		v1.AddFilter(atc.BEFORE_ROUTE,"users.*",AfterLogin)
//	}
//	v2 := atc.Route.Group("V2"){
//		v2.AddFilter(atc.BEFORE_ROUTE,"shop.car.*",AfterLogin)
//	}
func (group *RouterGroup) AddFilter(location Location, module string, filter FilterFunc) {
	moduleName := strings.Split(module, ".")
	module = path.Join(moduleName...)
	HttpAPP.Handler.AddFilter(location, path.Join("/", group.versionPath, module), filter)
}

// AddRouter add routing.
func AddRouter(module string, c HandlerInterface) {
	moduleName := strings.Split(module, ".")
	module = path.Join(moduleName...)
	HttpAPP.Handler.AddRouter(path.Join("/", module), c)
}

func AddFilter(location Location, module string, filter FilterFunc) {
	moduleName := strings.Split(module, ".")
	module = path.Join(moduleName...)
	HttpAPP.Handler.AddFilter(location, path.Join("/", module), filter)
}

func ExecuteHandler(httpMethod, module string, c *Handler) {
	HttpAPP.Handler.ExecuteHandler(httpMethod, path.Join("/", module), c)
}

// Shut down package method.

func NewGrace() *grace.Grace {
	if graceNodeTree == nil {
		graceNodeTree = grace.NewGrace()
	}
	return graceNodeTree
}