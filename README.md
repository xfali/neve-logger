# neve-logger

neve-logger是neve的日志扩展组件，用于配置日志。

内置日志中间件为[xlog](https://github.com/xfali/xlog)

## 安装
```
go get github.com/xfali/neve-logger
```

## 使用
  
### 1. neve集成（依赖[neve-core](https://github.com/xfali/neve-core)）
```
app := neve.NewFileConfigApplication("assets/config-test.yaml")
app.RegisterBean(xlogneve.NewLoggerProcessor())
//注册其他对象
// ...
app.Run()
```

### 2. 配置
在config-example.yaml中配置示例如下：
```
neve:
  logger:
    level: warn # [debug|info|warn|error|panic|fatal], default: info
    file:
      - ./neve-example.log
      - stdout
    # - stderr
    caller:
      file: short # [none|short|long], default: short
      func: simple # [none|short|simple|long], default: none
    simpleName: true
    #noFatalTrace: true
```
* 【neve.logger.file】配置日志输出文件，可以为系统日志文件路径以及：stdout表示输出到标准输出，stderr表示输出到标准错误输出。
* 【neve.logger.level】配置日志级别，包括debug、info、warn、error、panic、fatal
* 【neve.logger.caller.file】配置日志显示调用者文件格式，包括none、short、long
* 【neve.logger.caller.func】配置日志显示调用者方法格式，包括none、short、simple、long
* 【neve.logger.simpleName】配置是否xlog.GetLogger传入的对象使用简单名称而不使用全名
* 【neve.logger.noFatalTrace】配置Fatal日志输出时是否也输出堆栈信息

