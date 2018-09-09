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
	"testing"
)

func TestSetKey(t *testing.T) {
	key := []byte("123456789")
	SetKey(key)
	k := GetKey()
	if len(k) != 16 {
		t.Logf("%s\n", string(k))
		t.Fail()
	}

	if string(k) != "1234567890000000" {
		t.Logf("%s\n", string(k))
		t.Fail()
	}

	key = []byte("1234567812345678")
	SetKey(key)
	k = GetKey()

	if len(k) != 16 {
		t.Logf("%s\n", string(k))
		t.Fail()
	}

	if string(k) != "1234567812345678" {
		t.Logf("%s\n", string(k))
		t.Fail()
	}

}
