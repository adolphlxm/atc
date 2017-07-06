# logs

日志模块，目前支持的引擎有stdout、file

# 安装

    go get github.com/adolphlxm/atc/logs
   
# 使用步骤
##第一步：引入包
    
    import(
        "github.com/adolphlxm/atc/logs"
    )
    
##第二步：添加输出引擎

    logs.SetLogger("stdout")
    
    // 引擎支持第二个参数，配置信息
    logs.SetLogger("file",`{"filename":"`+AppConfig.DefaultString("log.file","")+`"}`)
    
## 通用方式(推进)

```go
package main

import (
    "github.com/adolphlxm/atc/logs"
)

func main() {

    logs.Debug("")
    logs.Info("")
    logs.Warn("")
    logs.Error("")
    ...
}

```

## 多实例使用

```go
   package main

    import (
        "github.com/adolphlxm/atc/logs"
    )

    func main() {
        log := logs.NewLogger(1024)
        log.SetLogger(logs.AdapterStdout)
        log.Debug("this is a debug message")
    }

```