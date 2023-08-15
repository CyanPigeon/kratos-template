# 鸽子群第六届字节青训营后端基础班CyanPigeon组 kratos框架模板代码
该代码根据[kratos官方的示例代码](https://github.com/go-kratos/kratos-layout)修改而来，主要是加了注释，加了TODO用来提示哪里写什么代码。

**注意**：最终合并到主分支的代码中**不要包含**任何的TODO，也**不要包含**`demo.go`和任何示例代码，`api/demo/v1`记得删了。

*原则上*要求按目录定义写。

*本地写的时候别用那个beta `Makefile`，你又没有make。到时候部署的时候再用。*

## 目录定义
```bash
toktik
└─ app
    ├─ cmd                 # 微服务主函数所在文件夹，一般不用动。
    ├─ configs             # 微服务配置文件，一般不用动。
    └─ internal            # 核心代码。
        ├─ biz             # 物理表定义与物理表操作的接口定义。MVC的M。
        │   └─ biz.go      # IoC相关，一般是来修改函数名的。对于biz，一般不需要管。
        ├─ conf            # 配置，一般也不用管。要改的话应该是大伙一起改。
        ├─ data            # 物理表操作的接口实现。MVC的M。
        │   └─ data.go     # IoC相关，一般是来修改函数名的。对于data，如果有多个物理表，需要在wire.NewSet()里面增加构造函数。如果没有物理表请留空。
        ├─ server          # 注册api的，一定有修改。
        │   └─ server.go   # IoC相关，一般是来修改函数名的。对于server，一般不需要管。
        └─ service         # 业务代码。MVC的C。
            └─ service.go  # IoC相关，一般是来修改函数名的。对于service，如果有多个api，需要在wire.NewSet()里面增加构造函数。
```
## 创建HTTP接口

> 参考 [kratos-gorm-git](https://github.com/getcharzp/kratos-gorm-git)
> 
> B站  [【项目实战】基于Kratos、gorm、git实现 代码托管平台](https://www.bilibili.com/video/BV17Y4y1y7jt)

> 在本标题 `创建HTTP接口` 下面提到的所有加 **{ }** 的表示 **自己填写补充的内容** ，**请自行替换**

首先创建一个proto文件

```bash
kratos proto add api/{package}/{api}.proto
```

如果加入版本号，则

```bash
kratos proto add api/{package}/{version}/{api}.proto
```

编辑proto文件，添加引用，删除多余的rpc和message，定义方法和路由

最后的文件内容将类似示例

> package => dy; 
> 
> api => feed;

```proto
syntax = "proto3";

package api.dy;

import "google/api/annotations.proto";
option go_package = "demo/api/dy;dy";

//定义一个叫Feed的service，其中FeedRequest为请求体，FeedReply为响应体
//option那一行必须写才能定义http接口，get表GET方法，后面跟路由
service Feed {
  rpc Feed (FeedRequest) returns (FeedReply){
    option(google.api.http) = {
      get:"/feed"
    };
  }
}

message FeedRequest {}
message FeedReply {}
```

codegen出来proto的go文件

```bash
kratos proto client api/{package}/{api}.proto
```

生成service

```bash
kratos proto server api/{package}/{api}.proto t internal/service
```

在 `internal/service/{api}.go` 生成了相应的service

打开 `internal/server/http.go` 文件

定位到 `v1.RegisterGreeterHTTPServer(srv, greeter)` 这一行，在下方加入

```go
	{package}.Register{API}HTTPServer(srv, service.New{API}Service())
```

启动服务端，调试接口

## 创建微服务
*很好，kratos的cli工具有猫饼，加了--nomod参数过后这个import跟喝了假酒一样*

*不过问题已经解决了，记得更新kratos的cli工具*
### 方案一
```bash
# 用的bash语法，记得把变量替换一下，免得创建出一堆乱七八糟的怪东西出来。

app_name=demo # app名，按需替换。
api_name=demo # api名，按需替换。一个微服务可以有多个api，一个个创建即可。
version=v1    # api版本，默认v1。

# 创建app。
kratos new app/${app_name} --nomod -r https://github.com/CyanPigeon/kratos-template.git

# 创建api的IDL模板. 要求: api的IDL路径为api/${api_name}/${version}/${api_name}.proto。
# IDL模板创建后，打开模板进行修改。
# 为了规范错误处理，要求必须创建error_reason.proto文件用于错误信息描述。error_reason.proto的内容见api/demo/v1/error_reason.proto。
kratos proto add api/${api_name}/${version}/${api_name}.proto

# 通过IDL文件生成API代码。
kratos proto client api/${api_name}/${version}/${api_name}.proto
kratos proto client api/${api_name}/v1/error_reason.proto

# 生成完API代码过后，直接开始写业务。业务代码完成后再进行下一步。

# 生成所有的proto代码、wire等。生成wire_gen.go的时候可能会失败，目前只观察到了依赖没有被使用导致的异常。照着模板写应该是不会失败的。
go generate ./...

# 运行微服务。
# 可能会提示有多个app，选要运行的那个app运行就行。
kratos run
```
由于某些不可抗的因素，`internal/server/http.go`和`internal/server/grpc.go`中的部分import会报错，删掉多余的路径就行。
### 方案二
直接在demo上面修改。这个方案应该不用多说啥。
## 服务发现
修改`internal/registrar/registrar.go#L22`为当前微服务的终结点的前缀，然后就可以使用服务发现功能。

举个例子，如果是用户登录的微服务，即实现的是[这一组接口](https://github.com/CyanPigeon/IDL/blob/main/user/user.proto)，那么上述位置的路径需要修改为`/douyin/user`。

**注意事项：**
1. 如果你的项目目录中没有`middleware/discovery`，请先使用`git pull`来更新代码。
2. 如果没有使用本仓库的模板来创建app，请参照`cmd/server/main.go`和`internal/registrar`中的内容自行修改自己的代码。
3. 如果是在本次提交之前就通过本仓库的模板创建了app，可以手动更新文件，只需要更新`cmd/server/main.go`和`internal/registrar`即可。
4. 务必修改`internal/registrar/registrar.go#L22`和`cmd/server/main.go#L25`，否则网关无法将请求路由至微服务。
## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

