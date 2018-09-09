package help

const Help = `			 gossh

NAME
	gossh is a smart ssh tool.It is developed by Go,compiled into a separate binary without any dependencies.

DESCRIPTION
		gossh can do the follow things:
		1.runs cmd on the remote host.
		2.push a local file or path to the remote host.
		3.pull remote host file to local.

USAGE
	1.Single Mode
		remote-comand:
		gossh -t cmd  -h host -P port(default 22) -u user(default root) -p passswrod [-f] command 

		Files-transfer:   
		<push file>   
		gossh -t push  -h host -P port(default 22) -u user(default root) -p passswrod [-f] localfile  remotepath 

		<pull file> 
		gossh -t pull -h host -P port(default 22) -u user(default root) -p passswrod [-f] remotefile localpath 

	2.Batch Mode
		Ssh-comand:
		gossh -t cmd -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] command 

		Files-transfer:   
		gossh -t push -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] localfile  remotepath 
		gosh -t pull -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] remotefile localpath

EMAIL
    	email.tata@qq.com 
`
