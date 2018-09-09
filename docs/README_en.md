# gossh 

## What's gossh?

gossh is a tool to  excute command on a remote machine or push a file or dir to remote machine. It can run single or batch mode .It use ssh protocol.

## Usage: 

you'd better to add gossh to the system $PATH.
ssh-comand:gossh -t cmd -h hostname -P port(default 22) -u username(default root) -p passswrod command 
ssh-scp:   gossh -t scp -h hostname -P port(default 22) -u username(default root) -p passswrod localfile  remotepath 

## Note:

1.if -h hostname is not given ,the gossh read ip.txt file from current directory to get ip ,then run the command or transe files. The contents of ip.txt is the ip list separated by newline .  

2.if -h hostname is given, only excute command on the -h hostname.

3.The host's password is get from password db automaticly,so register the password to the db before run command. 
  You can also use the '-p password' to given the password if you havn't registered the password to the db.

4.The default port is 22, you can change it by the '-P port' option.

5.The default user is root , you can change it by the '-u user' option.

## Warning:
1.You must make sure your command is safe ,because the comand is default running in batch mode.

2.Before you use 'gossh -t scp' to transfer files,make sure the remotepath is exist.It rewrite the exist file in the remote server default, make backup before run the command .

3.When you use 'gossh -t scp ' and not given the fullpath of localfile, it find the localfile an current path default.

3.You also can transfer the whole directory to  use 'gossh -t scp '. 

## Example:
1.show ip list of the server 10.238.48.101
gossh  -t cmd -h 10.238.48.101 "ip addr list"

2.show ip list of all server in the current path ip.txt file.
gossh -t cmd  "ip addr list"

3.Transfer the local file ip.txt of current path to the server 10.238.48.101:/data/tmp
gossh -t scp -h 10.238.48.101 ip.txt  /data/tmp  

4.Transfer the local file ip.txt of current path to all servers's /data/tmp  
gossh -t scp ip.txt  /data/tmp  

