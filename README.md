# 鸽子群第六届字节青训营后端基础班CyanPigeon组 kratos框架模板代码
该代码根据[kratos官方的示例代码](https://github.com/go-kratos/kratos-layout/tree/main)修改而来，主要是加了注释，加了TODO用来提示哪里写什么代码。

**注意**：最终合并到主分支的代码中**不要包含**任何的TODO，也**不要包含**`demo.go`和任何示例代码，`api/demo/v1`记得删了。

*原则上*要求按目录定义写。

*本地写的时候别用那个beta `Makefile`，你又没有make。到时候部署的时候再用。*

## 目录定义
```bash
toktik
└─ app
    ├─ cmd             # 微服务主函数所在文件夹，一般不用动。
    ├─ configs         # 微服务配置文件，一般不用动。
    └─ internal        # 核心代码。
        ├─ biz         # 物理表定义与物理表操作的接口定义。MVC的M。
        │   └─ biz.go  # IoC相关，一般是来修改函数名的。对于biz，一般不需要管。
        ├─ conf        # 配置，一般也不用管。要改的话应该是大伙一起改。
        ├─ data        # 物理表操作的接口实现。MVC的M。
        │   └─ data.go # IoC相关，一般是来修改函数名的。对于data，如果有多个物理表，需要在wire.NewSet()里面增加构造函数。如果没有物理表请留空。
        ├─ server      # 注册api的，一定有修改。
        │   └─ biz.go  # IoC相关，一般是来修改函数名的。对于server，一般不需要管。
        └─ service     # 业务代码。MVC的C。
            └─ biz.go  # IoC相关，一般是来修改函数名的。对于service，如果有多个api，需要在wire.NewSet()里面增加构造函数。
```

## 创建微服务
```bash
# 用的bash语法，记得把变量替换一下，免得创建出一堆乱七八糟的怪东西出来。

app_name=demo # app名，按需替换。
api_name=demo # api名，按需替换。一个微服务可以有多个api，一个个创建即可。
version=v1    # api版本，默认v1。

# 创建app。
kratos new ${app_name} --nomod -r https://github.com/CyanPigeon/kratos-template.git

# 创建api的IDL模板. 要求: api的IDL路径为api/${api_name}/${version}/${api_name}.proto。
# IDL模板创建后，打开模板进行修改。
# 为了规范错误处理，要求必须创建error_reason.proto文件用于错误信息描述。error_reason.proto的内容见api/demo/v1/error_reason.proto。
kratos proto add api/${api_name}/${version}/${api_name}.proto

# 通过IDL文件生成API代码。
kratos proto client api/${api_name}/${version}/${api_name}.proto
kratos proto client api/${api_name}/v1/error_reason.proto

# 通过IDL文件生成Controller层代码。
kratos proto server api/${api_name}/${version}/${api_name}.proto -t ${app_name}/internal/service

# 生成所有的proto代码、wire等。生成wire_gen.go的时候可能会失败，目前只观察到了依赖没有被使用导致的异常。照着模板写应该是不会失败的。
go generate ./...

# 运行微服务。
# 可能会提示有多个app，选要运行的那个app运行就行。
kratos run
```

## 路由注册
TODO，还没研究怎么把所有的微服务集中起来。

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

