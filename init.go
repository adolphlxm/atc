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
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/adolphlxm/atc/rpc/thrift"
)

var ThriftRPC thrift.Thrift

func init() {
	// Parsing configuration environment
	runMode := flag.String("m", "local", "Use -m <config mode>")
	configFile := flag.String("c", "./conf/app.ini", "use -c <config file>")
	version := flag.Bool("v", false, "Use -v <current version>")
	flag.Parse()

	// Show version
	if *version {
		fmt.Println("ATC version", VERSION, runtime.GOOS+"/"+runtime.GOARCH)
		fmt.Println("APP version", APPVERSION)
		os.Exit(0)
	}

	// Initialize app serve
	HttpAPP = NewApp()

	// 1. Initialize config
	initConfig(*configFile, *runMode)

	// 2. Initialize logs
	initLogs()

	// 3. Initalize cache
	if c := AppConfig.DefaultBool("cache.support", false); c {
		RunCaches()
	}

	// 4. Initalize queue publisher
	if Aconfig.QueuePublisherSupport {
		RunQueuePublisher()
	}

	// 5. Initalize queue consumer
	if Aconfig.QueueConsumerSupport {
		RunQueueConsumer()
	}

	// 6. Initalize mongodb
	if Aconfig.MgoSupport {
		RunMgoDBs()
	}

	// 7. Initalize orms
	if Aconfig.OrmSupport {
		RunOrms()
	}

	// 8. Initalize pgrpc
	if Aconfig.GrpcSupport {
		lazyInitGrpcServer()
	}
}

// Initalize thrift serve
func initThriftServe(){
	addr := Aconfig.ThriftAddr + ":" + Aconfig.ThriftPort
	ThriftRPC = thrift.NewThriftServe(`{"addr":"` + addr + `","secure":"` + Aconfig.ThriftSecure + `"}`)
	ThriftRPC.Debug(Aconfig.ThriftDebug)
	ThriftRPC.Factory(Aconfig.ThriftProtocol, Aconfig.ThriftTransport)
	ThriftRPC.Timeout(Aconfig.ThriftClientTimeout)
}