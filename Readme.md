
# bbrpc

BigBang blockchain (https://github.com/bigbangcore/BigBang) RPC client for golang! BigBang RPC >> bbrpc.

## **Notice**

使用这个库您需要自行承担潜在的bug带来的风险

## TODO
- [ ] more rpc method support
- [ ] websocket rpc support
- [ ] full rpc method test
- [ ] performance test
- [ ] full field match verify
- [ ] CI


## 测试
需要bigbang 在$PATH 下，即可以直接从命令行运行`bigbang`

基本测测试思路：go 启动bigbang，通过rpc交互，验证逻辑符合bibang规则

## 使用

```golang
client, err := bbrpc.NewClient(&bbrpc.ConnConfig{...})
tShouldNil(t, err, "failed to new rpc client")
defer client.Shutdown() //建议在系统关闭时调用shutdown
```

api 参考 https://github.com/bigbangcore/BigBang/wiki/JSON-RPC 

**Suggestion** Strong rpc password, EnableTLS better, Allowed ip from server side(bigbang) 

## Features

- majority of JSON-RPC API
- run bigbang server in golang (test usage)
- JSON-RPC data structures
- http.Client reuse,safe for concurrent use by multiple goroutines
- zero dependency
- TLS support

## version

主要的版本跟着BigBangCore 走，semantic version

建议的方式是选择与你的BigBang-server一致的版本，比如你的bigbang为 `0.9.1` 建议选择 `v0.9.x`

## QA

- http层如何复用http.Client的？
    - 每个bbrpc.Client 有一个http.Client，所有的请求发送到一个内部chan req (cap 100) ,内部有一个goroutine for select 这个chan req,一直处理请求直到shutdown，每个请求返回的结果写入到req 的chan resp 里。参考client.go

- JSON RPC request id 如何生成的？
    - `atomic.AddUint64(&c.id, 1)` 自增

- Why testing(`make test`) slow?
    - 目前暂时没有找到bigbang缩短出块周期的方法，有些测试需要等待出块，这个是主要的时间耗费
    - 暂时没有考虑`t.Parallel()`,这个需要对测试端口进行无冲突的管理分配
    - 启动bigbang后会等待1s

## LICENSE

TBD