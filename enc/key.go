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

package enc

import (
	"sync"
)

//aes default key
var (
	Key = []byte("suckdaNaanddf394des239")
	mu  = &sync.Mutex{}
)

//key长度为16,多了截取[:16]，少了补'0'
func SetKey(s []byte) {
	mu.Lock()
	defer mu.Unlock()
	n := len(s)
	if n < 16 {
		t := 16 - n
		for t > 0 {
			s = append(s, '0')
			t--
		}
	} else {
		s = s[:16]
	}
	//println(string(s))
	copy(Key, s)
}

//获取key
func GetKey() []byte {
	mu.Lock()
	defer mu.Unlock()
	return Key[:16]
}
