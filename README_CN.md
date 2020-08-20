# gossh

## 1.gossh是什么

gossh是一个使用go语言开发的极度简洁的ssh工具，只有一个二进制程序，无任何依赖，真正开箱即用。用于远程管理linux(类unix)机器：包括远程执行命令和推拉文件,支持单机和批量模式。

## 2.gossh能干什么

gossh提供3种核心功能：

1. 连接到远程主机执行命令。
2. 推送本地文件或者文件夹到远程主机。
3. 拉取远程主机的文件到本地。

![功能](https://github.com/andesli/gossh/raw/master/docs/images/gossh_function.png)

## 3.gossh运行模式

gossh支持单机模式和批量并行模式，也就是可以一次向一台机器发送命令执行，也可以一次向成千上万台台机器批量发送命令。批量并行模式也是gossh最大的一个特点，充分利用go在并发执行方面的优势。

1. 单机模式。
单机模式支持上文说的三种功能：远程执行命令，推送文件或者目录，拉取文件。

2. 批量模式。

可以通过-i 参数指定ip文件，通过-c 指定并发度。
批量并行模式同样支持上文说的三种功能：远程执行命令，推送文件或者目录，拉取文件。

### 并行和串行执行

1. 批量模式默认通过-c控制并发度，如果-c 设置为1默认是串行执行模式, -c 的值大于1是并行执行模式。
2. 并行执行模式下某台机器连不上或者执行失败不会自动退出，串行模式也一样，但是串行模式通过-s 参数可以使gossh执行过程中出错立即退出。

并行模式下没有提供出错退出的原因是，并行执行下，很难立即停止整个任务的执行，串行模式比较容易控制，在日常使用中，可以先使用串行模式验证功能，然后开启并行模式提升效率。


## 4.gossh用法

### 4.1程序获取

1.源码编译。

```
#需要有go编译环境
cd $GOPATH/src && git clone https://github.com/andesli/gossh.git
cd gossh

//gossh工具
go build ./cmd/gossh 

//密码加解密工具
go build ./cmd/passtool


//编译脚本编译amd64 386体系结构下windows和linux版本,放到./bin目录下，如果有其他体系结构需要使用也可以修改脚本执行编译。
./build.sh

```

2.如果不想从源码编译，编译好的二进制程序放bin/目录下。

得益于go语言优秀的跨平台特性，在./bin下已经为大家生成了amd64和386体系结构下windows和linux共计4个版本程序。

```
bin
|-- 386
|   |-- linux
|   |   |-- gossh
|   |   `-- passtool
|   `-- windows
|       |-- gossh.exe
|       `-- passtool.exe
`-- amd64
    |-- linux
    |   |-- gossh
    |   `-- passtool
    `-- windows
        |-- gossh.exe
        `-- passtool.exe
```

[点击立即下载](https://github.com/andesli/gossh/blob/master/bin)


### 4.2参数说明

- gossh

```
#gossh -h
flag needs an argument: -h
Usage of gossh:

  -t string
        running mode: cmd|push|pull (default "cmd")
        运行模式：cmd 远程执行命令，默认值；push 推送文件到远程； pull拉取远程文件到本地。
        
  -h string
        ssh ip
        
  -P string
        ssh port (default "22")
        ssh端口

  -u string
        ssh user (default "root")
        ssh用户名

  -p string
        ssh password
        密码
        

  -i string
        ip file when batch running mode (default "ip.txt")
        批量执行是指定ip文件，有关文件格式见下文。

  -c int
        the number of concurrency when b (default 30)
        批量并发执行的并发度，默认值是30，如果指定为1，gossh是串行执行。

  -s    if -s is setting, gossh will exit when error occurs
		-s是个bool型，只有到-c被指定为1时才有效，用来控制串行执行报错后是否立即退出。
        
  -e    password is Encrypted 
        如果密码传递的是密文，使用-e标记。-e适用于通过-p传递的密码和-i 指定的文件中存放的密码字段。

  -key string
        aes key for password decrypt and encryption
		密码加解密使用的key，gossh有一个默认的加密key,可以通过-key=xxx指定一个加解密的key. passtool密码加解密工具同样支持该-key选项.
        
  -f    force to run even if it is not safe
        如果遇到危险命令gossh默认是不执行，危险命令目前收录的有（"rm", "mkfs", "mkfs.ext3", "make.ext2", "make.ext4", "make2fs", "shutdown", "reboot", "init", "dd"）,可以通过-f强制执行，-f 是bool型参数，不指定默认为false。

  -s    if -s is setting, gossh will exit when error occurs
		如果-c=1，即并发度为1串行执行时，默认出错后会退出，使用-s标记不要退出，继续执行，在-c>1时，无论是否指定-s都不会出错退出。
        
  -l string
        log level (debug|info|warn|error (default "info")
        日志级别

  -logpath string
        logfile path (default "./log/")
		日志存放目录，默认是./log/
        
```
- passtool密码工具

```
./passtool -h
Usage of ./passtool:
  -d    指定密码密文生成明文
  -e    指定密码明文生成密文
  -key string
        AES加密密钥
```


### 4.3 批量运行时IP文件格式

-i ipfile 指定批量操作的ip文件,ipfile文件每行有4个字段ip|port|user|password，字段之间使用|分隔，四个字段分别是：机器IP，ssh端口，ssh用户名，ssh密码。其中ip字段是必须的，其他三个字段是选填的。
下面的配置都是合法的。

```
ip|port|user|password
ip|port|user|
ip|port|user
ip|port|
ip|port
ip|
ip
```
如果没有提供可选字段，gossh 默认通过-u -p -P参数从命令行参数获取，如果没有指定命令行参数，默认取命令行参数的默认值。
gossh 当前参数的默认值：

```
-u 默认值是root
-P 默认值是22
-p 默认值是空
-t 默认值是cmd

```
**说明**  

- 密码字段如果是空,gossh默认从db插件中查找相关流程参考第5章。
- 如果密码字段加密了，需要指定-e标记。-e是个整体开关：密码文件中的密码要么全部加密，要么不加密。

### 4.4 详细示例

点击[示例](https://github.com/andesli/gossh/blob/master/docs/example.md)查看详情。

### 4.5 输出和日志

点击[输出和日志](https://github.com/andesli/gossh/blob/master/docs/output_format.md)查看详情。


## 5.密码管理

点击[密码管理](https://github.com/andesli/gossh/blob/master/docs/password.md)查看详情。

## 6.安全性

gossh从多种角度保证执行安全，包括密码的加密存放、命令黑名单、以及文件传递过程中的检查、日志记录等，详情
点击[gossh安全管理](https://github.com/andesli/gossh/blob/master/docs/safe.md)查看详情。

## 7.不是重复造轮子

gossh不是重复制造一个像ansible的轮子，gossh的核心目标是提供给运维人员一个极度简洁的ssh工具，方便运维人员远程批量并行的初始化和管理机器。

有很多同学说ansible已经够好的了，为什么还要搞gossh?这是一个误区，请问ansible怎么批量安装到机器上？python环境怎么批量安装？这里有一个“先有鸡还是先有蛋”的问题，gossh就是第一个会下蛋的鸡。gossh使用go语言开发，静态编译为二进制程序，只要你的机器有ssh环境，并且能密码可以登录，理论上都能使用gossh进行管理。

gossh核心目标就是解决机器交付后“最初一公里-机器初始化的工作”。此时机器除了ssh，可能没有任何其他运行环境，此时通过gossh能够方便的快速的初始化机器。比如安装python,mysql等。

即便大公司在每个服务器上部署自研的agent,统一平台管理所有的服务器，但是也不能保证管理平台整个链路容灾高可用，gossh至少提供了一条备用的链路，能够在运维平台出问题的情况下，以闪电般的速度解决问题，从这个角度说，保留gossh这个最简单的通道，也算为运维人员提供一条消防通道。

当然gossh也提供了扩展，可以方便进行二次开发，将其改造为远程执行引擎，集成到公司的自动化系统中。

## 8.gossh适用场景。

1. 大规模机器的首次初始化。
公司来了几百台机器，只有ssh环境，除了初始用户名和密码，没有其他的安装。此时使用gossh对机器进行初始化，建立基本的环境。（gossh当初写的时候就是为了解决腾讯支付DB几千台机器的环境初始化）。

2. 命令行式批量远程管理。
不是每个公司都是BAT，建立起自动化的运维管理系统。占绝大多数的中小企业的运维人员是通过脚本在远程管理机器，他们迫切需要一个拿来就用，不需要任何依赖的ssh工具。gossh就是为这种人准备的，gossh不需要任何配置文件，没有任何依赖，真正做到拿来即用。

3. 将gossh二次开发，改造为一个远程执行引擎。目前gossh还没有实现，没有实现的原因有两点：

- 这和最初将gossh设计为极简的ssh工具背道而驰。
- 没有实际场景需求，大部分公司使用ansible或者salt stack, gossh无意重复造轮子。如果你的公司是一家使用go开发所有基础设施的云计算公司，可以在此基础上开发出一个远程执行引擎。

**gossh目标**

1.第一阶段目标是提供一个极简、好用、无任何依赖、可并发执行的ssh命令行工具。

2.第二阶段目标是是实现一个高度可集成的远程执行引擎，对外提供API服务，完善点对点的文件传递功能，很大概率不会在现有的gossh程序上改造，避免gossh变的臃肿，而是提供一个gossh2或者gosshweb的工具，专门做这件事情，遵照go的哲学，少即是多，一次只做一件事情。

现阶段之关注第一阶段目标，第二阶段目标还在酝酿中。

## 9. FAQ

[FAQ](https://github.com/andesli/gossh/blob/master/docs/faq.md)

任何问题可联系 <email.tata@qq.com>, 感觉有用的话帮忙加个星。

为方便大家使用，提供了一个qq技术群:851647540， 手机qq可以直接扫描下方二维码。
![qq群](https://github.com/andesli/gossh/raw/master/docs/images/gossh_qq.png)


