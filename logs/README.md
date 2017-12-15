# logs

日志模块，目前支持的引擎有stdout、file

# 安装

    go get github.com/adolphlxm/atc/logs
   
# 使用步骤
## 第一步：引入包
    
    import(
        "github.com/adolphlxm/atc/logs"
    )
    
## 第二步：添加输出引擎

    logs.SetLogger("stdout")
    
    // 引擎支持第二个参数，配置信息
    logFile := &File{
                      LogDir:"./",
                      MaxSize:1000,
                      Buffersize:1000,
                      FlushInterval:5,
                     }
    logs.SetLogger("file", logFile)

## 参数说明

* logdir 存放日志路径
* maxsize 日志分割文件尺寸, 单位:byte
* buffersize 日志缓冲区大小，单位:byte
* flushinterval 定时刷新日志到磁盘的间隔时间，单位:s

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
    // 程序退出时把缓冲区数据刷入磁盘
    logs.Flush()
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
        log.Flush()
    }

```