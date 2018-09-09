# gossh使用示例

## 1.单机模式

1.远程运行命令。

- 命令原型

```
gossh [-t cmd] -h hostname -P port(default 22) -u username(default root) -p passswrod  [-f]  "command"

```
- -t 参数可以省略，gossh默认执行模式就是远程执行命令。
- command远程执行命令如果有空格，建议使用双引号将括起来。
- -P -u 都可以不指定，使用默认值。

- 示例

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -t cmd -u root -p xxxx -P 22 "uname"
ip=192.168.56.2
command=uname
return=0
Linux

----------------------------------------------------------

```
- 使用默认值

只指定了-h 主机，-t -u -P都使用默认值。

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2  "uname"
ip= 192.168.56.2
command= uname
Linux
```

- 使用-f强制执行命令，一些危险命令拒绝执行，使用-f强制执行。

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2  "cd /tmp && rm ip.txt"
Dangerous command in cd /tmp && rm ip.txt
You can use the `-f` option to force to excute
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -f  "cd /tmp && rm ip.txt"
ip=192.168.56.2
command=cd /tmp && rm ip.txt
return=0

----------------------------------------------------------
```

2.推送文件到远程主机。

支持推送文件或者文件夹到远程主机，如果远程主机已经存在文件，默认是拒绝执行，可以使用-f参数强制覆盖。

- 命令原型

```
gossh -t push -h hostname -P port(default 22) -u username(default root) -p passswrod  [-f]  localfile/localpath  remotepath

```

如果远程文件已存在,可以使用-f参数强制覆盖。

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -t push ip.txt  /tmp
ip=192.168.56.2
command=push ip.txt to 192.168.56.2:/tmp
return=1
<ERROR>
Remote Server's /tmp has the same file ip.txt
You can use `-f` option force to cover the remote file.
</ERROR>

----------------------------------------------------------
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -t push -f ip.txt  /tmp
ip=192.168.56.2
command=push ip.txt to 192.168.56.2:/tmp
return=0
push ip.txt to 192.168.56.2:/tmp ok

----------------------------------------------------------
```

3.从远程主机拉取文件。

- 命令原型

```
gossh -t pull -h hostname -P port(default 22) -u username(default root) -p passswrod remote_file  local_path

```

注意：如果local_path不存在，会自动创建。

- 完整示例

```
[root@andesli.com /project/go/src/gossh]#gossh -h 192.168.56.2 -u root -p xxxxx -P 22 -t pull  /tmp/ip.txt .
ip= 192.168.56.2
command= scp  root@192.168.56.2:/tmp/ip.txt .
Files is transferred successfully.

```

## 2.批量模式

批量模式和单机模式的区别：

1. 使用-i参数代替-h指定批量运行的ip文件。
2. 可以是用-c参数指定并发度。
3. 如果密码加密了，需要带上-e参数告诉gossh密码已经加密。
4. 批量拉取文件到本地时，由于是文件都一样，文件会放到local_path/ip下，单机模式时直接放在local_path指定本地目录，但是批量模式，放到local_path下每个主机ip子目录下。

ipfile文件内容（密码加密）：

```
[root@andesli.com /project/go/src/gossh]#cat ip.txt 
192.168.56.2|22|root|T9GrQBSkD6zkRZOEd+ggfg==
192.168.56.2|22|root|
192.168.56.2|22
192.168.56.2|

```


1.批量远程运行命令。

```
[root@andesli.com /project/go/src/gossh]#gossh -t cmd -i ip.txt "uname"
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

2.推送文件到批量远程主机。

```
[root@andesli.com /project/go/src/gossh]#gossh -t push -f -i ip.txt  passtool /tmp   
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

3.批量从远程主机拉取文件到本地。

```
[root@andesli.com /project/go/src/gossh]#gossh -t pull -f -i ip.txt  /project/go/src/gossh/passtool /tmp 
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

注意：

-h hostname 指定主机，用于单机模式。  
-i ipfile 指定ip文件，用于批量模式。

如果既不指定-h 也不指定-i参数，gossh默认会从当前目录中寻找ip.txt文件作为ip文件进行批量模式执行。

文件推送或者拉取过程中，设计到文件是否存在的判断，以及是否覆盖的判断，详细规则见[gossh安全管理](https://github.com/andesli/gossh/blob/master/docs/safe.md)第3节文件传递安全性。

## 3 使用总结

gossh可以完成简单的命令执行，文件传递等工作，也可以完成复查的工作。完成复杂工作需要在本地编写脚本，推送脚本文件，然后远程执行脚本，再将脚本执行结果文件拉取到本地分析处理。

![用法](https://github.com/andesli/gossh/raw/master/docs/images/gossh_use.png)

