# gossh输出打印格式

## 标准输出 

1.gossh远程执行命令返回格式.

```
#批量模式首行首先打印所有的远程机器IP.
[servers]=[192.168.56.2 192.168.56.2]
#机器ip
ip=xxx.xxx.56.2
#远程执行命令
command=uname
#命令执行完后的退出值，就是$?
return=0
#远程执行命令输出到标准输出和错误输出的结果
Linux

##换行和---分隔线

----------------------------------------------------------
```

下面是一个简单的示例：

```
[root@andesli.com /project/go/src/gossh]#gossh "uname"
[servers]=[192.168.56.2 192.168.56.2]
ip=192.168.56.2
command=uname
return=0
Linux

----------------------------------------------------------
ip=192.168.56.2
command=uname
return=0
Linux

----------------------------------------------------------
```
远程执行命令在实际使用中的两个范式：

1. 如果是临时性的任务，一般gossh结合grep能够很方便的判断批量执行的结果。
2. 如果是正式任务或者复杂任务，建议将逻辑封装到一个脚本文件中，push到远程主机,然后再执行。

2.gossh推送和拉取文件输出结果和远程执行命令格式类似。

通过return=0判断推送或者拉取文件成功。

- push文件

```
[root@andesli.com /project/go/src/gossh]#gossh -t push passtool /tmp  
[servers]=[192.168.56.2 192.168.56.2]
ip=192.168.56.2
command=push passtool to 192.168.56.2:/tmp
return=0
push passtool to 192.168.56.2:/tmp ok

----------------------------------------------------------
ip=192.168.56.2
command=push passtool to 192.168.56.2:/tmp
return=0
push passtool to 192.168.56.2:/tmp ok

----------------------------------------------------------
```
- pull文件

```
[root@andesli.com /project/go/src/gossh]#gossh -t pull -f /project/go/src/gossh/passtool /tmp
[servers]=[192.168.56.2 192.168.56.2]
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/passtool /tmp/192.168.56.2
return=0
Pull from /project/go/src/gossh/passtool to /tmp/192.168.56.2 ok.
----------------------------------------------------------
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/passtool /tmp/192.168.56.2
return=0
Pull from /project/go/src/gossh/passtool to /tmp/192.168.56.2 ok.
----------------------------------------------------------

```


## 2.日志

1. gossh日志名默认是gossh.log ,该文件默认位于执行gossh当前目录的./log内(./log/gossh.log)，可以通过-logpath path 选项指定日志文件位置，如果目录不存在，gossh自动创建该目录，暂时不支持修改日志文件名。

2. 支持如下日志级别debug|info|warn|error, 默认是info级别，可以通过-l 选项指定日志级别。

3. gossh日志不会打印到标准输出里面，仅仅作为审计使用。

gossh日志模块使用的是beego的日志模块，采用异步记录模式，详见[https://github.com/astaxie/beego/tree/master/logs](https://github.com/astaxie/beego/tree/master/logs)，beego日志模块具有良好的扩展性，稍许定制就可以方便的将日志输出到文件，邮件和数据库中。



