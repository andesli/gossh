// Copyright 2018 github.com/andesli/gossh Author. All Rights Reserved.
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

package run

import (
	"fmt"
	"github.com/andesli/gossh/machine"
	"github.com/andesli/gossh/output"
	//	"net"
	//	"strings"
	//"context"
	"errors"
	"github.com/andesli/gossh/config"
	"github.com/andesli/gossh/logs"
	//	"github.com/andesli/gossh/tools"
	"path/filepath"
	"sync"
)

var (
	log = logs.NewLogger()
)

type CommonUser struct {
	user    string
	port    string
	psw     string
	force   bool
	encflag bool
}

func NewUser(user, port, psw string, force, encflag bool) *CommonUser {
	return &CommonUser{
		user:    user,
		port:    port,
		psw:     psw,
		force:   force,
		encflag: encflag,
	}

}
func SingleRun(host, cmd string, cu *CommonUser, force bool) {
	server := machine.NewCmdServer(host, cu.port, cu.user, cu.psw, "cmd", cmd, force)
	r := server.SRunCmd()
	output.Print(r)
}

//func ServersRun(cmd string, cu *CommonUser, wt *sync.WaitGroup, crs chan machine.Result, ipFile string, ccons chan struct{}) {
func ServersRun(cmd string, cu *CommonUser, wt *sync.WaitGroup, crs chan machine.Result, ipFile string, ccons chan struct{}, safe bool) {
	hosts, err := parseIpfile(ipFile, cu)
	if err != nil {
		log.Error("Parse %s error, error=%s", ipFile, err)
		return
	}

	ips := config.GetIps(hosts)

	//config.PrintHosts(hosts)
	log.Info("[servers]=%v", ips)
	fmt.Printf("[servers]=%v\n", ips)

	ls := len(hosts)

	//ccons==1 串行执行,可以暂停
	if cap(ccons) == 1 {
		log.Debug("串行执行")
		for _, h := range hosts {
			server := machine.NewCmdServer(h.Ip, h.Port, h.User, h.Psw, "cmd", cmd, cu.force)
			r := server.SRunCmd()
			if r.Err != nil && safe {
				log.Debug("%s执行出错", h.Ip)
				output.Print(r)
				break
			} else {
				output.Print(r)
			}
		}
	} else {
		log.Debug("并行执行")
		go output.PrintResults2(crs, ls, wt, ccons)

		for _, h := range hosts {
			ccons <- struct{}{}
			server := machine.NewCmdServer(h.Ip, h.Port, h.User, h.Psw, "cmd", cmd, cu.force)
			wt.Add(1)
			go server.PRunCmd(crs)
		}
	}
}

func SinglePush(ip, src, dst string, cu *CommonUser, f bool) {
	server := machine.NewScpServer(ip, cu.port, cu.user, cu.psw, "scp", src, dst, f)
	cmd := "push " + server.FileName + " to " + server.Ip + ":" + server.RemotePath

	rs := machine.Result{
		Ip:  server.Ip,
		Cmd: cmd,
	}
	err := server.RunScpDir()
	if err != nil {
		rs.Err = err
	} else {
		rs.Result = cmd + " ok\n"
	}
	output.Print(rs)
}

//push file or dir to remote servers
func ServersPush(src, dst string, cu *CommonUser, ipFile string, wt *sync.WaitGroup, ccons chan struct{}, crs chan machine.Result) {
	hosts, err := parseIpfile(ipFile, cu)
	if err != nil {
		log.Error("Parse %s error, error=%s", ipFile, err)
		return
	}

	ips := config.GetIps(hosts)
	log.Info("[servers]=%v", ips)
	fmt.Printf("[servers]=%v\n", ips)

	ls := len(hosts)
	go output.PrintResults2(crs, ls, wt, ccons)

	for _, h := range hosts {
		ccons <- struct{}{}
		server := machine.NewScpServer(h.Ip, h.Port, h.User, h.Psw, "scp", src, dst, cu.force)
		wt.Add(1)
		go server.PRunScp(crs)
	}
}
func SinglePull(host string, cu *CommonUser, src, dst string, force bool) {
	server := machine.NewPullServer(host, cu.port, cu.user, cu.psw, "scp", src, dst, force)
	err := server.PullScp()
	output.PrintPullResult(host, src, dst, err)
}

// pull romote server file to local
func ServersPull(src, dst string, cu *CommonUser, ipFile string, force bool) {
	hosts, err := parseIpfile(ipFile, cu)
	if err != nil {
		log.Error("Parse %s error, error=%s", ipFile, err)
		return
	}
	ips := config.GetIps(hosts)
	log.Info("[servers]=%v", ips)
	fmt.Printf("[servers]=%v\n", ips)

	for _, h := range hosts {
		ip := h.Ip

		localPath := filepath.Join(src, ip)
		server := machine.NewPullServer(h.Ip, h.Port, h.User, h.Psw, "scp", localPath, dst, cu.force)
		err = server.PullScp()
		output.PrintPullResult(ip, localPath, dst, err)
	}
}

//common logic
func parseIpfile(ipFile string, cu *CommonUser) ([]config.Host, error) {
	hosts, err := config.ParseIps(ipFile, cu.encflag)
	if err != nil {
		log.Error("Parse Ip File %s error,%s\n", ipFile, err)
		return hosts, err
	}

	if len(hosts) == 0 {
		return hosts, errors.New(ipFile + " is null")
	}
	hosts = config.PaddingHosts(hosts, cu.port, cu.user, cu.psw)
	return hosts, nil

}
