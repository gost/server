// Copyright (c) 2016 The GOST Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
)

var configLocation = "../config.yaml"
var configLocationDrone = "../config_drone.yaml"
var configFake = "nonexistingfile.yaml"
var configWrongData = "testreadconfig.yaml"

var data = `
server:
    name: GOST Server
    host: localhost
    port: 8080
    externalUri: localhost:8080/
database:
    host: 192.168.40.10
    port: 5432
    user: postgres
    password: postgres
    database: gost
    schema: v1
    ssl: false
mqtt:
    enabled: true
    webSocketEnabled: true
    secureWebSocketEnabled: false
    port: 1883
    webSocketPort:  8081
    secureWebSocketPort: 8082
    keepAlive: 300
    connectTimeout: 2
    ackTimeout: 20
    timeoutRetries: 3
    secureWebSocketCertPath:
    secureWebSocketKeyPath:
`

func TestReadFile(t *testing.T) {
	f, err := readFile(configLocation)
	if err != nil {
		t.Error("Please make sure there is a config.yaml file in the root directory ", err)
	}

	assert.NotNil(t, f, "config bytes should not be nil")

	f, err = readFile(configLocationDrone)
	if err != nil {
		t.Error("Error for config_drone.yaml", err)
	}

	assert.NotNil(t, f, "config bytes should not be nil")

	_, err = readFile(configFake)
	if err == nil {
		t.Error("Reading non existing config file should have given an error")
	}
}

func TestReadConfig(t *testing.T) {
	// try parsing data, this should parse if not give error
	content := []byte(data)
	cfg, err := readConfig(content)
	if err != nil {
		t.Error("Given static config data could not be parsed into config struct")
	}
	assert.NotNil(t, cfg)

	// check some random params
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, false, cfg.Database.SSL)
	assert.Equal(t, "192.168.40.10", cfg.Database.Host)

	// put in some false content this should fail
	falseContent := []byte("aaabbbccc")
	cfg, err = readConfig(falseContent)
	assert.NotNil(t, err, "ReadConfig should have returned an error")
}

func TestGetConfig(t *testing.T) {
	// Get the default config, should exist
	cfg, err := GetConfig(configLocation)
	if err != nil {
		t.Error("GetConfig returned an error ", err)
	}

	// Check if there is data
	assert.NotNil(t, cfg.Server.Host)
	assert.NotNil(t, cfg.Database.Database)

	// Try to get a fake config, should return error
	_, err = GetConfig(configFake)
	if err == nil {
		t.Error("Given fake config did not return an error", err)
	}

	// write file with some fake data, GetConfig should return error
	d1 := []byte("aaabbbccc")
	_ = ioutil.WriteFile(configWrongData, d1, 0644)
	_, err = GetConfig(configWrongData)
	_ = os.Remove(configWrongData)
	if err == nil {
		t.Error("GetConfig should have returned an error on reading a fake config", err)
	}
}
