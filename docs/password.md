# 密码管理

## 1.密码安全

gossh批量操作时，将密码明文存放到配置文件中有时不太妥。gossh提供加密存放方式。通过-e开关，默认是不加密存放，-e参数代表当前的IP配置文件是加密后存放的。

为此提供了一个专门的加解密工具passtool。

```
[root@andesli.com /project/go/src/gossh]#./passtool 
  -d    指定密码密文生成明文
  -e    指定密码明文生成密文
  -key string
        加密密钥
```
如果是长期运行的脚本，建议将ip中的密码加密处理。

## 2.密码获取

gossh对于密码的支持比较灵活，可以通过-p参数指定，批量模式下可以在ip文件中指定，ip文件中的的密码支持加密。如果上述都没有指定，gossh还默认的可以通过插件方式从db或者通过外部系统api获取密码。

- 单机模式密码获取流程

![单机模式](https://github.com/andesli/gossh/raw/v0.1/docs/images/singlepass.png)

- 批量模式模式密码获取流程

![批量模式](https://github.com/andesli/gossh/raw/v0.1/docs/images/batchpass.png)

## 3.密码扩展

gossh支持密码插件的方式访问密码，通过定义一套标准的密码获取接口，外部插件只要实现该接口，就能注册进去。gossh内部默认实现了一个通过db获取密码的插件。

1. 通过访问db获取密码。
gossh 提供了一个简单的默认实现，如果不指定操作的机器密码，gossh默认会访问：[gossh/auth/db/query.go](https://github.com/andesli/gossh/blob/master/auth/db/query.go#L9)中指定的db库表中查询。

```
 10 const (
 11     dbtype = "mysql" 
 13     ip       = "localhost"
 14     port     = "3306"
 15     user     = "mysql_user"
 16     passwd   = "mysql_pass"
 17     dbname   = "cmdb"

```
db库表初始化sql参见 [gossh/sql/db_init.sql](https://github.com/andesli/gossh/blob/master/sql/db_init.sql)

[gossh通过db获取密码环境搭建过程](https://github.com/andesli/gossh/blob/master/docs/use_mysql_db.md)

2. 通过web api方式。
该种方法只写了框架，需要的同学可以将其与自己的密码管理系统对接。

## 4.密码插件原理

gossh这个密码获取实现参照go标准库database/sql的设计思想，主要代码在auth目录。任何实现了如下接口的密码获取组件都可以注册到drivers里面。

- 密码接口,gossh/auth/driver/driver.go
```
  type GetPassworder interface {
       GetPassword(ip, user string) (string, error)
  }
```

- 注册接口.
```
// gossh/auth/auth.go

drivers   = make(map[string]driver.GetPassworder)
func Register(name string, d driver.GetPassworder)

```
- 密码存放到db的注册实现.

```
func init() {
    db := &DbDriver{
        dbtype:   dbtype,
        ipport:   ipport,
        user:     user,
        password: passwd,
        dbname:   dbname,
        sql:      querysql,
    }
    auth.Register("db", db)
}
```

- 注册或者修改密码获取组件。
```
// gossh/machine/server.go 

//加载注册实现了GetPassworder密码访问组件
  _ "gossh/auth/db"
  //_ "gossh/auth/web"
  
//指定从那个密码源读取密码
    PASSWORD_SOURCE = "db"
    //PASSWORD_SOURCE   = "web"

```
## 5. 加密key

加解密密码的默认key存放在 [key](https://github.com/andesli/gossh/blob/master/enc/key.go) ,gossh和passtool都支持通过-key选项指定加解密的key。


注意:  

1. -key指定的key要求有16个字节，如果key不够16个字节程序自动会在后面填充"0",如果超过16个字节，程序自动会截取16字节长度。
2. 目前存放在db里面的机器密码也是使用该key加解密。

## 6.SSH2登录认证

SSH2协议有三种登录方式，分别是：

```
Password
Keyboard Interaction
Public Key
```
gossh支持 Password 和Keyboard Interaction，理论上只要linux机器支持密码登录，都可以使用gossh进行管理，当然你也可以通过gossh为机器配置Public Key。

