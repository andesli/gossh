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

--create db user for gossh
create user 'mysql_user'@'localhost' identified by 'mysql_pass';
grant  select on cmdb.* to 'mysql_user'@'localhost';

--init server login information
--注意密码必须使用passtool加密,插入的是加密后的值
insert into t_password_info (hostName,userName,curPSW) values ('192.168.56.3','root','+ojuqnTp/hXWtEZSn5xE7w=='); 
insert into t_password_info (hostName,userName,curPSW) values ('192.168.56.2','root','+ojuqnTp/hXWtEZSn5xE7w=='); 

