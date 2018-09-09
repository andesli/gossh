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

//package main

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var (
	sourcekey  = []byte("zbOXeOdSedMK34QilHPUHw==")
	cleartext  = "tata"
	ciphertext = ""
)

func TestAesEncrypt(t *testing.T) {
	key := sourcekey[:16]
	result, err := AesEncrypt([]byte(cleartext), key)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("encryptData=", base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(origData) != cleartext {
		t.Fatal("dec error")
	}

}
