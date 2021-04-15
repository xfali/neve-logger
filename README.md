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
    level: warn
    file:
      - ./neve-example.log
      - stdout
```
* 【neve.logger.file】配置日志输出文件，可以为系统日志文件路径以及：stdout表示输出到标准输出，stderr表示输出到标准错误输出。
* 【neve.logger.level】配置日志级别，包括debug、info、warn、error、panic、fatal
