# gossh安全性

## 1.远程执行命令

### 1.1 危险命令检测.

gossh将危险的命令放到黑名单中，一旦远程执行危险命令，会自动退出，通过指定-f参数强制执行。危险命令目前收录如下：

```
"mount", "umount", "rm", "mkfs", "mkfs.ext3", "make.ext2", "make.ext4", "make2fs", "shutdown", "reboot", "init", "dd"
```

### 1.2 串行命令执行.

gossh为了提升效率，采用了并发执行模式，通过-c参数控制并发度，-c的默认值是30，当然也支持串行执行，-c=1就是串行执行，串行执行是如果遇到失败，gossh会自动退出，不会再继续执行，这也是一种安全防护，可以使用-s参数强制遇到执行失败仍然继续执行。非串行模式，gossh遇到一个host执行错误，不会退出，这是因为gossh内部是并发执行的，有一个出错，其他的并发执行的任务很难停下来。在实际使用gossh开发的过程中，建议先使用gossh -c=1 进行串行的调试，带安全测试后，再改为并行执行模式。

## 2.密码安全

gossh支持密码加密存放，如果指定了-e标记，就代表传递给gossh的密码密文。-e开关不但对单机执行有效，-i指定的ip文件里的存放的密码也被认为是加密存放的，密码的加解密可以使用passtool工具处理。

- gossh使用-key指定加解密key。

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -p="5f9lPu0eHz98lRsnpo+oHw==" -key="tata"   -e "uname"
ip=192.168.56.2
command=uname
return=0
Linux

----------------------------------------------------------
```
- passtool 指定key加解密

```
[root@andesli.com /project/go/src/gossh]#passtool -key="sos123" -e tata                       
a12AAcm9PaUmIppJvq7fFw==
[root@andesli.com /project/go/src/gossh]#passtool -key="sos123" -d a12AAcm9PaUmIppJvq7fFw==
tata

```

gossh也支持密码插件功能，可以很方便的将加密后的密码存放到数据库中，详情见[gossh通过db获取密码环境搭建过程](https://github.com/andesli/gossh/blob/master/docs/use_mysql_db.md)

## 3.文件传递安全性

### 3.1 push文件到远程

1.会检测本地文件是否存在，不存在报错。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t push  -f /project/go/src/gossh/ss /tmp
ip=192.168.56.2
command=push /project/go/src/gossh/ss to 192.168.56.2:/tmp
return=1
stat /project/go/src/gossh/ss: no such file or directory
----------------------------------------------------------
```
2.会检查远程目录是否存在，如果不存在报错。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t push  -f /project/go/src/gossh/gossh.go /tata
ip=192.168.56.2
command=push /project/go/src/gossh/gossh.go to 192.168.56.2:/tata
return=1
[192.168.56.2:/tata] does not exist or not a dir

----------------------------------------------------------
```
3.会检查远程目录是否存在同名文件，如果存在，默认报错，可以使用-f参数强制覆盖。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t push  /project/go/src/gossh/gossh.go /tmp  
ip=192.168.56.2
command=push /project/go/src/gossh/gossh.go to 192.168.56.2:/tmp
return=1
<ERROR>
Remote Server's /tmp has the same file /project/go/src/gossh/gossh.go
You can use `-f` option force to cover the remote file.
</ERROR>

----------------------------------------------------------
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t push  -f /project/go/src/gossh/gossh.go /tmp
ip=192.168.56.2
command=push /project/go/src/gossh/gossh.go to 192.168.56.2:/tmp
return=0
push /project/go/src/gossh/gossh.go to 192.168.56.2:/tmp ok

----------------------------------------------------------
```
### 3.2 pull文件安全性

1.pull会检测远程拉取的文件是否存在，不存在会报错。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t pull  /project/go/src/gossh/ss /tmp        
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/ss /tmp
return=1
Remote Server's /project/go/src/gossh/ss doesn't exist.

----------------------------------------------------------
```

2.pull会检测远程拉取的是否是文件，如果是目录会报错。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t pull  /project/go/src/gossh/ /tmp  
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/ /tmp
return=1
Remote Server's /project/go/src/gossh/ is a directory ,not support.

----------------------------------------------------------
```

3.pull会检查指定的本地路径是否存在，如果不存在会自动创建。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t pull  /project/go/src/gossh/gossh  /tmp/tata
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/gossh /tmp/tata
return=0
Pull from /project/go/src/gossh/gossh to /tmp/tata ok.
----------------------------------------------------------
```

4.pull会检查指定的本地路径是否存在，如果指定的不是目录而是文件也会报错。

```
[root@andesli.com /project/go/src/gossh]#./gossh -h 192.168.56.2 -t pull  /project/go/src/gossh/gossh  /tmp/gossh
ip= 192.168.56.2
command= scp  root@192.168.56.2:/project/go/src/gossh/gossh /tmp/gossh
return=1
/tmp/gossh is a normal file ,not a dir
----------------------------------------------------------
```
## 4. 日志审计

gossh 默认会将所有执行的命令记录到日志文件中，日志文件默认位于./log/gossh.log。可以通过-l 选项指定日志级别，可以通过-logpath指定日志位置。 具体点击[输出和日志](https://github.com/andesli/gossh/blob/master/docs/output_format.md)详细页面。


