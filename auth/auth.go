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

package auth

import (
	"errors"
	"gossh/auth/driver"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.GetPassworder)
)

//Register Password Source Driver
func Register(name string, d driver.GetPassworder) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if d == nil {
		panic("Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("Register called twice for driver " + name)
	}
	drivers[name] = d
}

//List Password Source Drivers
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]driver.GetPassworder)
}

// get password
func GetPassword(driverName, ip, user string) (string, error) {
	driversMu.Lock()
	defer driversMu.Unlock()
	d, ok := drivers[driverName]

	if !ok {
		return "", errors.New("unknown password driver: " + driverName)
	}
	return d.GetPassword(ip, user)
}
