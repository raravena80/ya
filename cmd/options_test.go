// Copyright © 2017 Ricardo Aravena <raravena80@gmail.com>
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
	"github.com/raravena80/ya/common"
	"github.com/spf13/viper"
	"testing"
)

func TestBuildCommonOptions(t *testing.T) {
	// Set up viper values for testing
	viper.Set("ya.machines", []string{"host1", "host2"})
	viper.Set("ya.user", "testuser")
	viper.Set("ya.port", 2222)
	viper.Set("ya.key", "/tmp/test.key")
	viper.Set("ya.useagent", true)
	viper.Set("ya.timeout", 30)

	options := BuildCommonOptions()

	// Create an Options struct and apply the options
	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// Verify all values were set correctly
	if len(opt.Machines) != 2 {
		t.Errorf("Expected 2 machines, got %d", len(opt.Machines))
	}
	if opt.Machines[0] != "host1" {
		t.Errorf("Expected host1, got %s", opt.Machines[0])
	}
	if opt.Machines[1] != "host2" {
		t.Errorf("Expected host2, got %s", opt.Machines[1])
	}
	if opt.User != "testuser" {
		t.Errorf("Expected testuser, got %s", opt.User)
	}
	if opt.Port != 2222 {
		t.Errorf("Expected port 2222, got %d", opt.Port)
	}
	if opt.Key != "/tmp/test.key" {
		t.Errorf("Expected /tmp/test.key, got %s", opt.Key)
	}
	if !opt.UseAgent {
		t.Error("Expected UseAgent to be true")
	}
	if opt.Timeout != 30 {
		t.Errorf("Expected timeout 30, got %d", opt.Timeout)
	}
}

func TestBuildCommonOptionsDefaults(t *testing.T) {
	// Reset viper to defaults
	viper.Reset()
	viper.Set("ya.machines", []string{})
	viper.Set("ya.user", "")
	viper.Set("ya.port", 0)
	viper.Set("ya.key", "")
	viper.Set("ya.useagent", false)
	viper.Set("ya.timeout", 0)

	options := BuildCommonOptions()

	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// Verify default values
	if opt.Machines == nil {
		t.Error("Expected Machines to be initialized, got nil")
	}
	if opt.User != "" {
		t.Errorf("Expected empty user, got %s", opt.User)
	}
	if opt.Port != 0 {
		t.Errorf("Expected port 0, got %d", opt.Port)
	}
	if opt.UseAgent {
		t.Error("Expected UseAgent to be false")
	}
}
