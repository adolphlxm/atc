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
	"path/filepath"

	"github.com/adolphlxm/atc/logs"
	"github.com/adolphlxm/atc/rpc/thrift"
)

var ThriftRPC thrift.Thrift

func init() {
	// Parsing configuration environment
	runMode := flag.String("m", "dev", "Use -m <config mode>")
	configFile := flag.String("c","../conf/app.ini","use -c <config file>")
	version := flag.Bool("v", false, "Use -v <current version>")
	flag.Parse()

	// Show version
	if *version {
		fmt.Println("ATC version",VERSION, runtime.GOOS +"/" + runtime.GOARCH)
		os.Exit(0)
	}

	// Initialize config
	err := ParseConfig(*configFile,*runMode)
	if err != nil {
		workPath, _ := os.Getwd()
		workPath, _ = filepath.Abs(workPath)
		fmt.Printf("workPath: %v", workPath)
		panic(err)
	}

	logs.SetLogger(Aconfig.LogOutput, `{"filename":"`+AppConfig.DefaultString("log.file","")+`", "perm":"`+AppConfig.DefaultString("log.perm","")+`"}`)
	if Aconfig.Debug {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(Aconfig.LogLevel)
	}

	// Initializes error
	ErrorCode = NewErrorMap()
	// In the conf/error. Ini file parsing error code
	err = ErrorCode.parse(AppConfig.DefaultString("error.file","../conf/error.ini"))
	if err != nil {
		logs.Error("Error file loading err:%v", err.Error())
	}

	// Initialize app serve
	HttpAPP = NewApp()

	// Initalize thrift serve
	addr := Aconfig.ThriftAddr + ":" + Aconfig.ThriftPort
	ThriftRPC = thrift.NewThriftServe(`{"addr":"` + addr + `","secure":"` + Aconfig.ThriftSecure + `"}`)
	ThriftRPC.Debug(Aconfig.ThriftDebug)
	ThriftRPC.Factory(Aconfig.ThriftProtocol, Aconfig.ThriftTransport)
	ThriftRPC.Timeout(Aconfig.ThriftClientTimeout)
}
