# gossh通过db获取密码环境搭建过程

gossh对于密码的获取非常灵活，除了最简单的从命令行参数和ip文件中指定外，gossh在内部还实现了一个插件，那就是将密码初始化到DB中，gossh自动从db获取密码进行执行。现在就来介绍下如何搭建这样一个环境。

## 1.搭建一个mysql环境。

这里不再详解，参见[install mysql](https://dev.mysql.com/doc/refman/5.7/en/installing.html).

## 2.初始化库表。

```
--create db and table
create database if not exists cmdb; use cmdb;
create table if not exists  t_password_info (
		  hostName varchar(225) NOT NULL DEFAULT '' COMMENT '登录机器ip',
		  userName varchar(225) NOT NULL DEFAULT ''COMMENT '登录机器用户名',
		  curPSW varchar(225) NOT NULL DEFAULT '' COMMENT '登录机器当前密码',
		  curPSWStatus int(11) NOT NULL DEFAULT 0 COMMENT '当前密码状态',
		  expiredTime datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '密码过期时间',
		  lastPSW1 varchar(225) NOT NULL DEFAULT '',
		  lastPSW2 varchar(225) NOT NULL DEFAULT '',
		  createTime datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
		  modifyTime datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
		  defaultPSW varchar(225) NOT NULL DEFAULT 'tmpPasword',
		  flag int(11) NOT NULL DEFAULT '0',
		  PRIMARY KEY (hostName,userName)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

```

## 3.创建gossh访问mysql的用户。

如果mysql不在gossh所在的机器，将localhost改为gossh部署的机器ip。

```
--create db user for gossh
create user 'mysql_user'@'localhost' identified by 'mysql_pass';
grant  select on cmdb.* to 'mysql_user'@'localhost';

```
## 4.初始化机器登录信息到db.

将要管理的机器元信息初始化到DB库表中，注意密码字段要使用passtool加密后的值。

```
--init server login information
insert into t_password_info (hostName,userName,curPSW) values ('localhost','root','+ojuqnTp/hXWtEZSn5xE7w=='); 
insert into t_password_info (hostName,userName,curPSW) values ('192.168.56.1','root','+ojuqnTp/hXWtEZSn5xE7w=='); 

```

## 5.修改gossh中的db配置参数

gossh代码中写死了连接db的配置（这一块未来考虑通过参数可以指定），如果访问db的IP、端口、用户名、密码和代码中不一致，需要根据实际情况修改,代码位置如下：

```
// gossh/auth/db/query.go

  9 const (   
 10     dbtype   = "mysql"
 11     ipport   = "localhost:3306"
 12     user     = "mysql_user"
 13     passwd   = "mysql_pass"
 14     dbname   = "cmdb"
 15     querysql = `select curPSW from t_password_info as A where A.hostName=? and A.userName= ? `
 16 ) 

```

## 6.编译一个二进制程序。

重新编译gossh程序。

```
#需要有go编译环境
cd $GOPATH/src && git clone https://github.com/andesli/gossh.git
cd gossh

//gossh工具
go build gossh.go
```
至此gossh就可以通过访问db获取密码了，这只是一个简单的实现，密码表设计的很简单。密码表每个公司设计各不一样，可以根据实际情况做改造，也可以自己实现一个密码插件注册进去，将gossh接入特定的密码管理系统中。

