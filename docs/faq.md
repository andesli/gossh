#  常见问题

1.gossh是否支持拉取远程目录？
目前gossh还不支持拉取远程机器的目录到本地，建议先将远程机器目录压缩成文件，然后在拉取。

2.gossh远程命令中的引号有什么好的建议？
执行命令comand如果较复杂，中间有空格，或者引号，使用原则是：

- 使用双引号将command整个引用起来 "command" ，此时command中如包含双引号需要使用\进行转义(\")，含有的单引号不需要任何处理，直接使用；

- 不建议使用单引号将 command扩其来 ，原因是单引号中不能再引用单引号（引号就近匹配原则，且其内部引用的东西不做任何转义，导致command书写的灵活性降低大大降低，特别是在脚本中。）

示例：

```
gossh -t cmd "ps -ef |egrep \"(mysql|master|slave|time|keep)\"" 
gossh -t cmd "ps -ef |egrep '(mysql|master|slave|time|keep)'" 

```
3.执行报"GET PASSWORD ERROR"原因。

gossh优先从ip文件中获取密码，然后从命令行获取密码，最后试图从db获取密码。如果你没有配置DB相关环境，也没有在命令行和IP文件中指定密码，gossh获取密码失败会报该错误。

4.gossh运行平台。

gossh使用go语言编写，理论上只要go语言支持的平台都能使用gossh，但是由于gossh使用的是ssh2协议，被gossh管理的机器必须支持ssh2协议，gossh最适合管理linux系统，windows系统没有经过测试。


5.任何问题请联系 email.tata@qq.com

为方便大家使用，提供了一个qq技术群:851647540， 可直接扫描下方二维码。
![qq群](https://github.com/andesli/gossh/raw/master/docs/images/gossh_qq.png)
