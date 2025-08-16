// Copyright © 2017 Ricardo Aravena <raravena@branch.io>
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

package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigWithFile(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
ya:
  user: testuser
  port: 2222
  key: /tmp/testkey
  useagent: true
  timeout: 10
  machines:
    - 192.168.1.1
    - 192.168.1.2
`)
	tmpfile, err := ioutil.TempFile("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a new viper instance
	v := viper.New()
	v.SetConfigFile(tmpfile.Name())
	v.SetConfigType("yaml")

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err != nil {
		t.Errorf("Failed to read config file: %v", err)
	}

	// Check if the values are loaded correctly
	if v.GetString("ya.user") != "testuser" {
		t.Errorf("Expected user to be 'testuser', but got '%s'", v.GetString("ya.user"))
	}
	if v.GetInt("ya.port") != 2222 {
		t.Errorf("Expected port to be 2222, but got '%d'", v.GetInt("ya.port"))
	}
	if v.GetString("ya.key") != "/tmp/testkey" {
		t.Errorf("Expected key to be '/tmp/testkey', but got '%s'", v.GetString("ya.key"))
	}
	if !v.GetBool("ya.useagent") {
		t.Errorf("Expected useagent to be true, but got false")
	}
	if v.GetInt("ya.timeout") != 10 {
		t.Errorf("Expected timeout to be 10, but got '%d'", v.GetInt("ya.timeout"))
	}
	machines := v.GetStringSlice("ya.machines")
	if len(machines) != 2 {
		t.Errorf("Expected 2 machines, but got %d", len(machines))
	}
	if machines[0] != "192.168.1.1" {
		t.Errorf("Expected machine 1 to be '192.168.1.1', but got '%s'", machines[0])
	}
	if machines[1] != "192.168.1.2" {
		t.Errorf("Expected machine 2 to be '192.168.1.2', but got '%s'", machines[1])
	}
}
