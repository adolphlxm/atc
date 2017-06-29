# 设计架构思路

## 整体设计架构
![图片](https://github.com/adolphlxm/atc/blob/dev/doc/image/module.png)

#请求生命周期

## http 

### HTTP Handler
    
    Hander？，它是一个接口。这个接口很简单，只要某个struct有ServeHTTP(http.ResponseWriter, *http.Request)这个方法，那这个struct就自动实现了Hander接口
    
* ATC 构建自 GO HTTP server， 他为每个请求创建一个 goroutine(轻量级线程)，用于并发处理。
* ATC 把 request 请求通过Handler交给 过滤器、Actions处理，完成后把结果写到response响应中。
    
        Handler接口中有一个 ServeHTTP(ResponseWriter, *Request) 方法,当SERVER从TCP端口中获取到新的请求时,会调用这个方法去执行具体,换而言之, ServeHTTP方法就是所有SERVER消息处理接口的入口函数.其他所有想要处理HTTP请求的方法都必须直接或间接通过这个接口实现.
    
### 过滤器
支持三种过滤器：
1. EFORE_ROUTE                    //匹配路由之前
2. BEFORE_HANDLER                 //匹配到路由后,执行Handler actions之前
3. AFTER                          //执行完actions逻辑后

### RESTFul
RESTful 是一种目前 API 开发中广泛采用的形式，ATC 支持这样的请求方法，也就是用户 Get 请求就执行 Get 方法，Post 请求就执行 Post 方法。

### Response
失败响应：
```go
{
  "error": "Undefined error",
  "code": 10003,
  "request": "/V1/ram/auth/"
}
```
成功响应：
```go
{
  "a": "1",
  "b": 2,
  "c": "3"
}
```
## RPC
### thrift
