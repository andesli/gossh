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

package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/andesli/gossh/auth"
	"github.com/andesli/gossh/enc"
)

const (
	dbtype   = "mysql"
	ip       = "localhost"
	port     = "3306"
	user     = "mysql_user"
	passwd   = "mysql_pass"
	dbname   = "cmdb"
	querysql = `select curPSW from t_password_info as A where A.hostName=? and A.userName= ? `
)

type DbDriver struct {
	dbtype   string
	ip       string
	port     string
	user     string
	password string
	dbname   string
	sql      string
}

func init() {
	db := &DbDriver{
		dbtype:   dbtype,
		ip:       ip,
		port:     port,
		user:     user,
		password: passwd,
		dbname:   dbname,
		sql:      querysql,
	}
	auth.Register("db", db)
}

func (dv DbDriver) GetPassword(host, user string) (string, error) {

	dbhost := "tcp(" + dv.ip + ":" + dv.port + ")"
	conn := dv.user + ":" + dv.password + "@" + dbhost + "/" + dv.dbname + "?" + ""

	db, err := sql.Open(dv.dbtype, conn)
	if err != nil {
		return "", err
	}
	if err = db.Ping(); err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(dv.sql)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	curPsw := ""
	err = stmt.QueryRow(host, user).Scan(&curPsw)

	if err != nil {
		return "", err
	}
	skey := enc.GetKey()

	psw, err := enc.AesDecEncode(curPsw, skey[:16])

	return string(psw), err

	//maybe need decrypt the password
	//return Decrypt(curPsw,key)
}
