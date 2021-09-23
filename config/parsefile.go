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

package config

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/andesli/gossh/enc"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

//Password Encrypted Key

var (
	rzkey = enc.GetKey()
)

type Host struct {
	Ip   string
	Port string
	User string
	Psw  string
}

/*
func main() {
	ipfile := os.Args[1]
	iplist, _ := ParseIps(ipfile, true)
	fmt.Printf("%v\n", iplist)
}
*/
func GetIps(h []Host) []string {
	ips := make([]string, 0)
	for _, v := range h {
		ips = append(ips, v.Ip)
	}
	return ips
}

func PrintHosts(h []Host) {
	for _, v := range h {
		fmt.Printf("host=%s,port=%s,user=%s,password=%s\n", v.Ip, v.Port, v.User, v.Psw)
	}
}

func PaddingHosts(h []Host, port, user, psw string) []Host {

	hosts := make([]Host, 0)
	for _, v := range h {
		if v.Port == "" {
			v.Port = port
		}
		if v.User == "" {
			v.User = user
		}
		if v.Psw == "" {
			v.Psw = psw
		}

		hosts = append(hosts, v)
	}
	return hosts
}

/*
ParseIps parse ip file to []Host

ip file supported format:
ip
ip|port
ip|port|username
ip|port|username|password
*/

func ParseIps(ipfile string, eflag bool) ([]Host, error) {
	hosts := make([]Host, 0)

	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	configfile := ""

	if filepath.IsAbs(ipfile) {
		configfile = ipfile
	} else {
		configfile = filepath.Join(AppPath, ipfile)
	}

	f, err := os.Open(configfile)

	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)

	for {
		s, err := buf.ReadString('\n')
		if err != nil {
			//主要是兼容windows和linux文件格式
			if err == io.EOF && s != "" {
				goto Lable
			} else {
				return hosts, nil
			}
		}
	Lable:
		line := strings.TrimSpace(s)
		//if line[0] == '#' || line == "" {
		if line == "" || line[0] == '#' {
			continue
		}
		h, err := parseLine(line, eflag)
		if err != nil {
			continue
			//	return hosts, err
		}
		hosts = append(hosts, h)
	}
}

func parseLine(s string, eflag bool) (Host, error) {
	host := Host{}
	line := strings.TrimSpace(s)

	if line[0] == '#' {
		return host, errors.New("comment line")
	}
	if line == "" {
		return host, errors.New("null line")
	}

	fields := strings.Split(line, "|")
	//ip := net.ParseIP(fields[0])
	hname := strings.TrimSpace(fields[0])
	_, err := net.LookupHost(hname)
	if err != nil {
		return host, errors.New("ill ip")
	}
	lens := len(fields)
	switch lens {
	case 1:
		host.Ip = hname
	case 2:
		host.Ip = hname
		host.Port = strings.TrimSpace(fields[1])
	case 3:
		host.Ip = hname
		host.Port = strings.TrimSpace(fields[1])
		host.User = strings.TrimSpace(fields[2])
	case 4:
		host.Ip = hname
		host.Port = strings.TrimSpace(fields[1])
		host.User = strings.TrimSpace(fields[2])
		pass := strings.TrimSpace(fields[3])
		if eflag && pass != "" {
			text, err := decrypt(pass, rzkey)
			if err != nil {
				return host, errors.New("decrypt the password error")
			}
			host.Psw = string(text)
		} else {
			host.Psw = pass
		}

	default:
		return host, errors.New("format err")
	}
	return host, nil
}

//decrypt password feild
func decrypt(pass string, key []byte) ([]byte, error) {
	skey := key[:16]
	return enc.AesDecEncode(pass, skey)
}
