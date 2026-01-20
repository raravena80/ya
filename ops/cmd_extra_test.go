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

package ops

import (
	"fmt"
	"testing"

	"github.com/raravena80/ya/common"
	"golang.org/x/crypto/ssh"
)

func TestExecuteCmd_ErrorPaths(t *testing.T) {
	// Test executeCmd with invalid hostname (will fail at Dial)
	tests := []struct {
		name     string
		hostname string
		options  common.Options
	}{
		{name: "Invalid hostname",
			hostname: "invalid.host.local",
			options: common.Options{
				Port: 22,
				Cmd:  "ls",
			}},
		{name: "Invalid port in options",
			hostname: "localhost",
			options: common.Options{
				Port: -1, // Invalid port
				Cmd:  "ls",
			}},
		{name: "Empty command",
			hostname: "localhost",
			options: common.Options{
				Port: 22,
				Cmd:  "",
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ssh.ClientConfig{
				User: "testuser",
				Auth: []ssh.AuthMethod{ssh.PublicKeys(nil)},
			}

			result := executeCmd(tt.options, tt.hostname, config)

			// Should return an error (connection will fail)
			if result.err == nil {
				t.Error("Expected error for invalid connection, got nil")
			}

			// Result should contain hostname
			if result.result == "" {
				t.Error("Expected result to contain hostname")
			}
		})
	}
}

func TestMakeExecResult_AllPaths(t *testing.T) {
	tests := []struct {
		name    string
		hostname string
		output  string
		err     error
	}{
		{name: "Success case",
			hostname: "testhost",
			output: "command output",
			err: nil},
		{name: "Error case",
			hostname: "testhost",
			output: "",
			err: fmt.Errorf("connection failed")},
		{name: "Output with error",
			hostname: "testhost",
			output: "partial output",
			err: fmt.Errorf("partial failure")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := makeExecResult(tt.hostname, tt.output, tt.err)

			// Check result format
			expectedPrefix := tt.hostname + ":\n"
			if tt.output != "" {
				expectedPrefix = tt.hostname + ":\n" + tt.output
			}

			if result.result != expectedPrefix {
				t.Errorf("Result = %q, want %q", result.result, expectedPrefix)
			}

			// Check error is preserved
			if result.err != tt.err {
				t.Errorf("Error = %v, want %v", result.err, tt.err)
			}
		})
	}
}
