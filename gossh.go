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

package main

import (
	"flag"
	"fmt"
	"gossh/enc"
	"gossh/help"
	"gossh/logs"
	"gossh/machine"
	"gossh/run"
	"gossh/tools"
	"path/filepath"
	"strings"
	"sync"
)

//gossh version
const (
	AppVersion = "gossh 0.6"
)

var (

	//common options
	port     = flag.String("P", "22", "ssh port")
	host     = flag.String("h", "", "ssh ip")
	user     = flag.String("u", "root", "ssh user")
	psw      = flag.String("p", "", "ssh password")
	prunType = flag.String("t", "cmd", "running mode: cmd|push|pull")

	//batch running options
	ipFile = flag.String("i", "ip.txt", "ip file when batch running mode")
	cons   = flag.Int("c", 30, "the number of concurrency when b")

	//safe options
	encFlag   = flag.Bool("e", false, "password is Encrypted")
	force     = flag.Bool("f", false, "force to run even if it is not safe")
	psafe     = flag.Bool("s", false, "if -s is setting, gossh will exit when error occurs")
	pkey      = flag.String("key", "", "aes key for password decrypt and encryption")
	blackList = []string{"rm", "mkfs", "mkfs.ext3", "make.ext2", "make.ext4", "make2fs", "shutdown", "reboot", "init", "dd"}

	//log options
	plogLevel = flag.String("l", "info", "log level (debug|info|warn|error")
	plogPath  = flag.String("logpath", "./log/", "logfile path")
	log       = logs.NewLogger()
	logFile   = "gossh.log"

	pversion = flag.Bool("version", false, "gossh version")
)

//main
func main() {

	usage := func() {
		fmt.Println(help.Help)
	}

	flag.Parse()

	//version
	if *pversion {
		fmt.Println(AppVersion)
		return
	}

	if *pkey != "" {
		enc.SetKey([]byte(*pkey))
	}

	if flag.NArg() < 1 {
		usage()
		return
	}

	if *prunType == "" || flag.Arg(0) == "" {
		usage()
		return
	}

	if err := initLog(); err != nil {
		fmt.Printf("init log error:%s\n", err)
		return
	}

	//异步日志，需要最后刷新和关闭
	defer func() {
		log.Flush()
		log.Close()
	}()

	log.Debug("parse flag ok , init log setting ok.")

	switch *prunType {
	//run command on remote server
	case "cmd":
		if flag.NArg() != 1 {
			usage()
			return
		}

		cmd := flag.Arg(0)

		if flag := tools.CheckSafe(cmd, blackList); !flag && *force == false {
			fmt.Printf("Dangerous command in %s", cmd)
			fmt.Printf("You can use the `-f` option to force to excute")
			log.Error("Dangerous command in %s", cmd)
			return
		}

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		log.Info("gossh -t=cmd  cmd=[%s]", cmd)

		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SingleRun(*host, cmd, puser, *force)

		} else {
			cr := make(chan machine.Result)
			ccons := make(chan struct{}, *cons)
			wg := &sync.WaitGroup{}
			run.ServersRun(cmd, puser, wg, cr, *ipFile, ccons, *psafe)
			wg.Wait()
		}

	//push file or dir  to remote server
	case "scp", "push":

		if flag.NArg() != 2 {
			usage()
			return
		}

		src := flag.Arg(0)
		dst := flag.Arg(1)
		log.Info("gossh -t=push local-file=%s, remote-path=%s", src, dst)

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SinglePush(*host, src, dst, puser, *force)
		} else {
			cr := make(chan machine.Result, 20)
			ccons := make(chan struct{}, *cons)
			wg := &sync.WaitGroup{}
			run.ServersPush(src, dst, puser, *ipFile, wg, ccons, cr)
			wg.Wait()
		}

	//pull file from remote server
	case "pull":
		if flag.NArg() != 2 {
			usage()
			return
		}

		//本地目录
		src := flag.Arg(1)
		//远程文件
		dst := flag.Arg(0)
		log.Info("gossh -t=pull remote-file=%s  local-path=%s", dst, src)

		puser := run.NewUser(*user, *port, *psw, *force, *encFlag)
		if *host != "" {
			log.Info("[servers]=%s", *host)
			run.SinglePull(*host, puser, src, dst, *force)
		} else {
			run.ServersPull(src, dst, puser, *ipFile, *force)
		}

	default:
		usage()
	}
}

//setting log
func initLog() error {
	switch *plogLevel {
	case "debug":
		log.SetLevel(logs.LevelDebug)
	case "error":
		log.SetLevel(logs.LevelError)
	case "info":
		log.SetLevel(logs.LevelInfo)
	case "warn":
		log.SetLevel(logs.LevelWarn)
	default:
		log.SetLevel(logs.LevelInfo)
	}

	logpath := *plogPath
	err := tools.MakePath(logpath)
	if err != nil {
		return err
	}

	logname := filepath.Join(logpath, logFile)
	logstring := `{"filename":"` + logname + `"}`

	//此处主要是处理windows下文件路径问题,不做转义，日志模块会报如下错误
	//logs.BeeLogger.SetLogger: invalid character 'g' in string escape code
	logstring = strings.Replace(logstring, `\`, `\\`, -1)

	err = log.SetLogger("file", logstring)
	if err != nil {
		return err
	}
	//开启日志异步提升性能
	log.Async()
	return nil
}
