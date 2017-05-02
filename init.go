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
	"github.com/adolphlxm/atc/logs"
	"github.com/adolphlxm/atc/rpc/thrift"
)

var ThriftRPC thrift.Thrift

func init() {
	// Initialize config
	err := ParseConfig()
	if err != nil {
		panic(err)
	}

	// Initialize log
	Logger = logs.NewLogger(10000)
	Logger.SetHandler(Aconfig.LogOutput, `{"filename":"`+AppConfig.DefaultString("log.file","")+`"}`)
	Logger.SetLevel(Aconfig.LogLevel)

	// Initialize app serve
	HttpAPP = NewApp()

	// Initalize thrift serve
	addr := Aconfig.ThriftAddr + ":" + Aconfig.ThriftPort
	ThriftRPC = thrift.NewThriftServe(`{"addr":"` + addr + `","secure":"` + Aconfig.ThriftSecure + `"}`)
	ThriftRPC.Debug(Aconfig.ThriftDebug)
	ThriftRPC.Factory(Aconfig.ThriftProtocol, Aconfig.ThriftTransport)
	ThriftRPC.Timeout(Aconfig.ThriftClientTimeout)
}
