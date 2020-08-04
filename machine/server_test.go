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
	"testing"
)

func initServer() *Server {
	s := &Server{
		Ip:     "127.0.0.1",
		Port:   "22",
		User:   "root",
		Action: "cmd",
		Cmd:    "uname",
		// password
		Psw: "hello@123",
	}

	return s
}

/*
func TestSetPsw(t *testing.T) {
	s := initServer()
	s.SetPsw()
	if s.Psw != "NO_PASSWORD" {
		t.Error("get password fail")
	}
}
*/

func TestRunCmd(t *testing.T) {
	s := initServer()
	//	s.SetPsw()

	_, err := s.RunCmd()
	if err != nil {
		t.Fail()
	}
}
