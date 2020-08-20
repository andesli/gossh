# gossh

[中文](https://github.com/andesli/gossh/blob/master//README_CN.md)

## 1.What's gossh

gossh is an extremely concise ssh tool which developed by go language. It has only a binary program without any dependencies and is really ready to use out of the box.
gossh is used Used to manage of linux (like unix) machines: including remote execution of commands and push and pull files, and support stand-alone and batch modes.


## 2.What can gossh do

Three core functions of gossh:

1. Connect to the remote host to execute commands.
2. Push local files or folders to remote hosts.
3. Pull files from the remote host to the local.

![功能](https://github.com/andesli/gossh/raw/master/docs/images/gossh_function.png)

## 3.gossh operating mode

gossh supports stand-alone mode and batch parallel mode, that is, it can send commands to one machine at a time for execution, or batch commands to thousands of machines at a time. The batch parallel mode is also one of the biggest features of gossh, making full use of the advantages of the go language in concurrent execution.
**Stand-alone mode**:  
The stand-alone mode supports the three functions mentioned above: remotely execute commands, push files or directories, and pull files.

**Batch mode**:

The ip file can be specified by the -i parameter, and the concurrency can be specified by -c parameter. The batch parallel mode also supports the three functions mentioned above: remotely execute commands, push files or directories, and pull files.

###  Execution mode :concurrent and serial 

1. The batch mode is controlled by -c by default. If -c is set to 1, the default is serial execution mode, and the value of -c is greater than 1 is parallel execution mode. 
2. In parallel execution mode, a machine cannot be connected or execution fails and will not automatically exit. Serial mode is the same, but serial mode can make gossh exit immediately when an error occurs during execution through the -s parameter. 

The reason why the error exit is not provided in the parallel mode is that it is difficult to stop the execution of the entire task immediately under the parallel execution. The serial mode is easier to control. In daily use, you can use the serial mode verification function first, and then turn on the parallel mode to improve effectiveness.


## 4.Getting started

### 4.1Install

**1.Building from source**

```
#To build gossh from the source code yourself you need to have a working Go environment with version 1.12 or greater installed.

cd $GOPATH/src && git clone https://github.com/andesli/gossh.git
cd gossh

//build gossh
go build ./cmd/gossh 

//build password encryption and decryption tool 
go build ./cmd/passtool


//Compile the programs for windows and linux os under the amd64 386 architecture, which  binarys is under the ./bin directory
./build.sh

```

**2.Pre-compiled binary**


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

[dowload](https://github.com/andesli/gossh/blob/master/bin)


### 4.2Usage

- gossh

```
#gossh -h
flag needs an argument: -h
Usage of gossh:

  -t string
        running mode: cmd|push|pull (default "cmd")
        
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

  -c int
        the number of concurrency when b (default 30)

  -s    if -s is setting, gossh will exit when error occurs
        
  -e    password is Encrypted 

  -key string
        aes key for password decrypt and encryption
        
  -f    force to run even if it is not safe

  -s    if -s is setting, gossh will exit when error occurs
        
  -l string
        log level (debug|info|warn|error (default "info")

  -logpath string
        logfile path (default "./log/")
        
```
- passtool tool

```
./passtool -h
Usage of ./passtool:
  -d    Convert ciphertext to plaintext
  -e    Convert plaintext to ciphertext
  -key string
        AES key
```


### 4.3 Config file 

The -i parameter is used to specify the batch operation host ip file. Each line of the file has 4 fields ip|port|user|password, separated by |. The four fields are: machine IP, ssh port, ssh user name, ssh password. The ip field is required, and the other three fields are optional. The following configurations are all legal.

```
ip|port|user|password
ip|port|user|
ip|port|user
ip|port|
ip|port
ip|
ip
```
If no optional fields are provided, gossh obtains the command line parameters through the -u, -p, -P parameters by default. If no command line parameters are specified, the default values of the command line parameters are taken by default. The default value of the current parameters of gossh:

```
-u root
-P 22
-p default empty 
-t cmd

```
**Remark**  

- If the password field is empty, gossh will find the relevant process from the db plugin by default, refer to 5.
- If the password field is encrypted, you need to specify the -e flag. -e is an overall switch: the passwords in the password file are either all encrypted or not.

### 4.4 Example 

[example](https://github.com/andesli/gossh/blob/master/docs/example.md)detail。

### 4.5 Log 

[logs](https://github.com/andesli/gossh/blob/master/docs/output_format.md)detail。


## 5.Password management

[Password management](https://github.com/andesli/gossh/blob/master/docs/password.md)detail。

## 6.Security

[Safety management](https://github.com/andesli/gossh/blob/master/docs/safe.md)detail


## 8.Scenes

1.The first initialization of a large-scale machine.

The company came to hundreds of machines, only the ssh environment, except the initial user name and password, no other installation. At this time, use gossh to initialize the machine and establish a basic environment. (When gossh was originally written, it was to solve the environment initialization of Tencent pay DB thousands of machines).

2.Command-line batch remote management.

Not every company is a BAT and has established an automated operation and maintenance management system. The operation and maintenance personnel of the vast majority of small and medium-sized enterprises manage machines remotely through scripts. They urgently need an ssh tool that can be used without any dependency. Gossh is prepared for this kind of people. Gossh does not require any configuration files, does not have any dependencies, and is really ready to use.


## 9. FAQ

[FAQ](https://github.com/andesli/gossh/blob/master/docs/faq.md)

Contact me for any questions<email.tata@qq.com>



