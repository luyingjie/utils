# utils
基础包
版本： v2.0.1


# 功能
    类型转换: convert

## v2.1.0
    添加网络部分
    添加Session
    添加协程池
    命令行解析
    配置文件全局读取

## v2.0.1
    重写error部分，报错会返回更多信息。
    错误会返回堆栈信息。
    编码解码
    转换
    万能变量
    json处理
    时间类型管理

## v2.0.0
    调整并简化文件结构。



## 计划完成
    注释调整
     有些组件没有写日志
    log还不够完善，信息比较冗余。
    文件处理
    出去基础包直接报错的代码
    
    基础类型
    http请求封装改善
    锁
    线程池
    命令行
    定时
    任务调度
    数据校验

    fsnotify中的异常处理还有问题。


### net/proxy问题

 - proxy代理consoel和boss的api请求的时候有错误，console表现在签名错误，替换成proxy2是ok的。
 - proxy2代理前端资源烨都正常，但是处理302的时候收到是200 ok。
 - proxy代理云管是ok的，都正常，但是使用自己处理返回值的时候有问题，这里应该不是这麽处理的。
 - proxy代理console这些前端资源的时候烨有问题。
    