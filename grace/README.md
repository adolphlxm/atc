# grace

顺序平滑退出包，使用双向链表list实现
退出顺序：从队尾开始逆序退出

# 安装

    go get github.com/adolphlxm/atc/grace
   
# 使用步骤
## 第一步：引入包
    
    import(
        "github.com/adolphlxm/atc/grace"
    )
    
## 第二步：实现接口

客户端退出需实现如下接口(支持多模块退出)：

    type TT interface {
    	ModuleID() string
    	Stop() error
    }

 * ModuleID() 方法返回string, 表示退出模块名称
 * Stop() 方法根据客户端业务自行实现退出逻辑


## 第三步：初始化Grace

    grace := grace.NewGrace()
    // 双向链表队头插入退出接口
    grace.PushFront(TT)
    // 双向链表队尾插入退出接口
    grace.PushBack(TT)
    // 在"atc"模块的链表之后插入退出接口
    grace.InsertAfter("atc",TT)
    // 在"atc"模块的链表之前插入退出接口
    grace.InsertBefore("atc",TT)
    // 在链表中移除"atc"退出接口
    grace.Remove("atc")
    // 在链表中将"atc"接口移动到"http"之后。
    grace.MoveAfter("atc", "http")
    // 在链表中将"atc"接口移动到"http"之前。
    grace.MoveAfter("atc", "http")
    
## 案例

```go
package main

import (
    "github.com/adolphlxm/atc/grace"
)

type ShutDown struct{}
func (this *ShutDown) ModuleID() string {
    return "down"
}
func (this *ShutDown) Stop() error {
    fmt.Println("客户端退出逻辑，完成退出")
    return nil
}

func main() {

    grace := grace.NewGrace()
    grace.PushFront(&ShutDown{})
    // ...
}

```