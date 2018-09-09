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

package hex

func HexStringToBytes(hex string) []byte {
	len := len(hex) / 2
	result := make([]byte, len)
	i := 0
	for i = 0; i < len; i++ {
		pos := i * 2
		result[i] = ToByte(hex[pos])<<4 | ToByte(hex[pos+1])
	}
	return result
}

func ToByte(c uint8) byte {

	if c >= '0' && c <= '9' {
		return byte(c - '0')
	}
	if c >= 'a' && c <= 'z' {
		return byte(c - 'a' + 10)
	}
	if c >= 'A' && c <= 'Z' {
		return byte(c - 'A' + 10)
	}
	return 0
}

func BytesToHexString(data []byte) string {
	hex := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	nLen := len(data)
	buff := make([]byte, 2*nLen)
	for i := 0; i < nLen; i++ {
		buff[2*i] = hex[(data[i]>>4)&0x0f]
		buff[2*i+1] = hex[data[i]&0x0f]
	}
	szHex := string(buff)
	return szHex
}
