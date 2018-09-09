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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	testAes()
}

func testAes() {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	sourcekey := []byte("zbOXeOdSedMK34QilHPUHw==")
	//key := ZeroPadding(sourcekey, 16)
	//key := []byte("sfe023f_9fd&fwfl")
	key := sourcekey[:16]
	result, err := AesEncrypt([]byte("3XZalfo*JV"), key)
	if err != nil {
		panic(err)
	}
	fmt.Println("encryptData=", base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("origData=", string(origData))
}

func AesEncEncode(origData, key []byte) (string, error) {

	si, err := AesEncrypt(origData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(si), nil
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	//origData = PKCS5Padding(origData, blockSize)
	origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecEncode(encodeStr string, key []byte) ([]byte, error) {
	crypted, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(crypted, key)
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	//origData = PKCS5UnPadding(origData)
	//fmt.Println("len=", len(origData))
	//fmt.Println("origData=", string(origData))
	//fmt.Println("origData2=", origData)
	origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	//fmt.Println("length-befor=", length)
	for length > 0 {
		unpadding := int(origData[length-1])
		if unpadding == 0 && length > 0 {
			length = length - 1
		} else {
			break
		}
	}
	//fmt.Println("length-after=", length)
	//return origData[:(length - unpadding)]
	if length == 0 {
		return origData[:1]
	} else {
		return origData[:length]
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
