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

package config

import (
	"testing"
)

func TestGetIps(t *testing.T) {
	hosts := make([]Host, 0)
	host1 := Host{
		Ip:   "192.168.56.2",
		Port: "22",
		User: "root",
		Psw:  "root",
	}
	host2 := Host{
		Ip:   "192.168.56.2",
		Port: "22",
		User: "root",
		Psw:  "root",
	}
	hosts = append(hosts, host1, host2)

	h := GetIps(hosts)
	if len(h) != 2 {
		t.Fatal("TestGetIps Error")
	}

}
