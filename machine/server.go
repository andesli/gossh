// Copyright 2018 gossh Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Author: andes
// Email: email.tata@qq.com

package machine

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"gossh/auth"
	_ "gossh/auth/db"
	//_ "gossh/auth/web"
	"gossh/logs"
	"gossh/scp"
	"gossh/tools"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	PASSWORD_SOURCE = "db"
	//PASSWORD_SOURCE   = "web"

	NO_PASSWORD = "GET PASSWORD ERROR\n"

	log = logs.NewLogger()
)

const (
	NO_EXIST = "0"
	IS_FILE  = "1"
	IS_DIR   = "2"
)

type Server struct {
	Ip         string
	Port       string
	User       string
	Psw        string
	Action     string
	Cmd        string
	FileName   string
	RemotePath string
	Force      bool
}

type ScpConfig struct {
	Src string
	Dst string
}

type Result struct {
	Ip     string
	Cmd    string
	Result string
	Err    error
}

func NewCmdServer(ip, port, user, psw, action, cmd string, force bool) *Server {
	server := &Server{
		Ip:     ip,
		Port:   port,
		User:   user,
		Action: action,
		Cmd:    cmd,
		Psw:    psw,
		Force:  force,
	}
	if psw == "" {
		server.SetPsw()
		//log.Debug("server.Psw=%s", server.Psw)
	}
	return server
}

func NewScpServer(ip, port, user, psw, action, file, rpath string, force bool) *Server {
	rfile := path.Join(rpath, path.Base(file))
	cmd := createShell(rfile)
	server := &Server{
		Ip:         ip,
		Port:       port,
		User:       user,
		Psw:        psw,
		Action:     action,
		FileName:   file,
		RemotePath: rpath,
		Cmd:        cmd,
		Force:      force,
	}
	if psw == "" {
		server.SetPsw()
	}
	return server
}
func NewPullServer(ip, port, user, psw, action, file, rpath string, force bool) *Server {
	cmd := createShell(rpath)
	server := &Server{
		Ip:         ip,
		Port:       port,
		User:       user,
		Psw:        psw,
		Action:     action,
		FileName:   file,
		RemotePath: rpath,
		Cmd:        cmd,
		Force:      force,
	}
	if psw == "" {
		server.SetPsw()
	}
	return server
}

/*
func NewScp(src, dst string) ScpConfig {
	scp := ScpConfig{
		Src: src,
		Dst: dst,
	}
	return scp
}
*/

//query password from password plugin
//PASSWORD_SOURCE: db|web
func (server *Server) SetPsw() {
	psw, err := auth.GetPassword(PASSWORD_SOURCE, server.Ip, server.User)
	if err != nil {
		server.Psw = NO_PASSWORD
		return
	}
	server.Psw = psw
}

//run command for parallel
func (server *Server) PRunCmd(crs chan Result) {
	rs := server.SRunCmd()
	crs <- rs
}

// set Server.Cmd
func (s *Server) SetCmd(cmd string) {
	s.Cmd = cmd
}

//run command in sequence
func (server *Server) RunCmd() (result string, err error) {
	if server.Psw == NO_PASSWORD {
		return NO_PASSWORD, nil
	}
	client, err := server.getSshClient()
	if err != nil {
		return "getSSHClient error", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "newSession error", err
	}
	defer session.Close()

	cmd := server.Cmd
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(bs), err
	}
	return string(bs), nil
}

//run command in sequence
func (server *Server) SRunCmd() Result {
	rs := Result{
		Ip:  server.Ip,
		Cmd: server.Cmd,
	}

	if server.Psw == NO_PASSWORD {
		rs.Err = errors.New(NO_PASSWORD)
		return rs
	}

	client, err := server.getSshClient()
	if err != nil {
		rs.Err = err
		return rs
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		rs.Err = err
		return rs
	}
	defer session.Close()

	cmd := server.Cmd
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		rs.Err = err
		return rs
	}
	rs.Result = string(bs)
	return rs
}

//execute a single command on remote server
func (server *Server) checkRemoteFile() (result string) {
	re, _ := server.RunCmd()
	return re
}

//PRunScp() can transport  file or path to remote host
func (server *Server) PRunScp(crs chan Result) {
	cmd := "push " + server.FileName + " to " + server.Ip + ":" + server.RemotePath
	rs := Result{
		Ip:  server.Ip,
		Cmd: cmd,
	}
	result := server.RunScpDir()
	if result != nil {
		rs.Err = result
	} else {
		rs.Result = cmd + " ok\n"
	}
	crs <- rs
}

func (server *Server) RunScpDir() (err error) {
	re := strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//远程机器存在同名文件
	if re == IS_FILE && server.Force == false {
		errString := "<ERROR>\nRemote Server's " + server.RemotePath + " has the same file " + server.FileName + "\nYou can use `-f` option force to cover the remote file.\n</ERROR>\n"
		return errors.New(errString)
	}

	rfile := server.RemotePath
	cmd := createShell(rfile)
	server.SetCmd(cmd)
	re = strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//远程目录不存在
	if re != IS_DIR {
		errString := "[" + server.Ip + ":" + server.RemotePath + "] does not exist or not a dir\n"
		return errors.New(errString)
	}

	client, err := server.getSshClient()
	if err != nil {
		return err
	}
	defer client.Close()

	filename := server.FileName
	fi, err := os.Stat(filename)
	if err != nil {
		log.Debug("open source file %s error\n", filename)
		return err
	}
	scp := scp.NewScp(client)
	if fi.IsDir() {
		err = scp.PushDir(filename, server.RemotePath)
		return err
	}
	err = scp.PushFile(filename, server.RemotePath)
	return err
}

//pull file from remote to local server
func (server *Server) PullScp() (err error) {

	//判断远程源文件情况
	re := strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//不存在报错
	if re == NO_EXIST {
		errString := "Remote Server's " + server.RemotePath + " doesn't exist.\n"
		return errors.New(errString)
	}

	//不支持拉取目录
	if re == IS_DIR {
		errString := "Remote Server's " + server.RemotePath + " is a directory ,not support.\n"
		return errors.New(errString)
	}

	//仅仅支持普通文件
	if re != IS_FILE {
		errString := "Get info from Remote Server's " + server.RemotePath + " error.\n"
		return errors.New(errString)
	}

	//本地目录
	dst := server.FileName
	//远程文件
	src := server.RemotePath

	log.Debug("src=%s", src)
	log.Debug("dst=%s", dst)

	//本地路径不存在，自动创建
	err = tools.MakePath(dst)
	if err != nil {
		return err
	}

	//检查本地是否有同名文件
	fileName := filepath.Base(src)
	localFile := filepath.Join(dst, fileName)

	flag := tools.FileExists(localFile)
	log.Debug("flag=%v", flag)
	log.Debug("localFile=%s", localFile)

	//-f 可以强制覆盖
	if flag && !server.Force {
		return errors.New(localFile + " is exist, use -f to cover the old file")
	}

	//执行pull
	client, err := server.getSshClient()
	if err != nil {
		return err
	}
	defer client.Close()

	scp := scp.NewScp(client)
	err = scp.PullFile(dst, src)
	return err
}

//RunScp1() only can transport  file to remote host
func (server *Server) RunScpFile() (result string, err error) {
	client, err := server.getSshClient()
	if err != nil {
		return "GetSSHClient Error\n", err
	}
	defer client.Close()

	filename := server.FileName
	session, err := client.NewSession()
	if err != nil {
		return "Create SSHSession Error", err
	}
	defer session.Close()

	go func() {
		Buf := make([]byte, 1024)
		w, _ := session.StdinPipe()
		defer w.Close()
		//File, err := os.Open(filepath.Abs(filename))
		File, err := os.Open(filename)
		if err != nil {
			log.Debug("open scp source file %s error\n", filename)
			return
		}
		defer File.Close()

		info, _ := File.Stat()
		newname := filepath.Base(filename)
		fmt.Fprintln(w, "C0644", info.Size(), newname)
		for {
			n, err := File.Read(Buf)
			fmt.Fprint(w, string(Buf[:n]))
			if err != nil {
				if err == io.EOF {
					// transfer end with \x00
					fmt.Fprint(w, "\x00")
					return
				} else {
					fmt.Println("read scp source file error")
					return
				}
			}
		}
	}()

	cmd := "/usr/bin/scp -qt " + server.RemotePath
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(bs), err
	}
	return string(bs), nil
}

// implement ssh auth method [password keyboard-interactive] and [password]
func (server *Server) getSshClient() (client *ssh.Client, err error) {
	authMethods := []ssh.AuthMethod{}
	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {

		if len(questions) == 0 {
			return []string{}, nil
		}
		/*
			for i, question := range questions {
				log.Debug("SSH Question %d: %s", i+1, question)
			}
		*/

		answers = make([]string, len(questions))
		for i := range questions {
			yes, _ := regexp.MatchString("*yes*", questions[i])
			if yes {
				answers[i] = "yes"

			} else {
				answers[i] = server.Psw
			}
		}
		return answers, nil
	}
	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(server.Psw))

	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: authMethods,
	}
	//psw := []ssh.AuthMethod{ssh.Password(server.Psw)}
	//Conf := ssh.ClientConfig{User: server.User, Auth: psw}
	ip_port := server.Ip + ":" + server.Port
	client, err = ssh.Dial("tcp", ip_port, sshConfig)
	return
}

//create shell script for running on remote server
func createShell(file string) string {
	s1 := "bash << EOF \n"
	s2 := "if [[ -f " + file + " ]];then \n"
	s3 := "echo '1'\n"
	s4 := "elif [[ -d " + file + " ]];then \n"
	s5 := `echo "2"
else 
echo "0"
fi
EOF`
	cmd := s1 + s2 + s3 + s4 + s5
	return cmd
}
