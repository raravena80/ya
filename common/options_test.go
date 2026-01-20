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

package common

import (
	"testing"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name      string
		machines  []string
		port      int
		timeout   int
		user      string
		cmd       string
		key       string
		src       string
		dst       string
		agentSock string
		op        string
		useAgent  bool
		recursive bool
		verbose   bool
	}{
		{name: "Test all options ssh",
			machines:  []string{"one", "two", "three"},
			port:      22,
			user:      "bogus",
			cmd:       "runit",
			key:       "mykey",
			src:       "src",
			dst:       "dst",
			timeout:   20,
			agentSock: "socket",
			op:        "run",
			useAgent:  false,
			recursive: true,
			verbose:   false,
		},
		{name: "Test all options scp",
			machines:  []string{"one", "two", "three"},
			port:      22,
			user:      "bogus",
			cmd:       "runit",
			key:       "mykey",
			src:       "src",
			dst:       "dst",
			timeout:   20,
			agentSock: "socket",
			op:        "copy",
			useAgent:  false,
			recursive: false,
			verbose:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []func(*Options){SetMachines(tt.machines),
				SetUser(tt.user),
				SetPort(tt.port),
				SetCmd(tt.cmd),
				SetKey(tt.key),
				SetSource(tt.src),
				SetDestination(tt.dst),
				SetTimeout(tt.timeout),
				SetAgentSock(tt.agentSock),
				SetAgentSock(tt.op),
				SetUseAgent(tt.useAgent),
				SetIsRecursive(tt.recursive),
				SetVerbose(tt.verbose),
				SetOp(tt.op)}
			opt := Options{}
			for _, option := range options {
				option(&opt)
			}
		})
	}
}

func TestSetKnownHosts(t *testing.T) {
	tests := []struct {
		name     string
		knownHosts string
		expected string
	}{
		{name: "Set known_hosts path",
			knownHosts: "/home/user/.ssh/known_hosts",
			expected: "/home/user/.ssh/known_hosts"},
		{name: "Set empty known_hosts",
			knownHosts: "",
			expected: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetKnownHosts(tt.knownHosts)(&opt)
			if opt.KnownHosts != tt.expected {
				t.Errorf("SetKnownHosts() = %v, want %v", opt.KnownHosts, tt.expected)
			}
		})
	}
}

func TestSetInsecureHost(t *testing.T) {
	tests := []struct {
		name          string
		insecureHost  bool
		expected      bool
	}{
		{name: "Set insecure host true",
			insecureHost: true,
			expected: true},
		{name: "Set insecure host false",
			insecureHost: false,
			expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetInsecureHost(tt.insecureHost)(&opt)
			if opt.InsecureHost != tt.expected {
				t.Errorf("SetInsecureHost() = %v, want %v", opt.InsecureHost, tt.expected)
			}
		})
	}
}

func TestOptionsStructFields(t *testing.T) {
	// Verify that all option setters actually set the correct fields
	opt := Options{}

	SetKnownHosts("/test/known_hosts")(&opt)
	if opt.KnownHosts != "/test/known_hosts" {
		t.Errorf("KnownHosts not set correctly, got: %v", opt.KnownHosts)
	}

	SetInsecureHost(true)(&opt)
	if !opt.InsecureHost {
		t.Error("InsecureHost not set correctly")
	}

	// Also verify the option function works as expected
	opt2 := Options{}
	optionFunc := SetInsecureHost(true)
	optionFunc(&opt2)
	if !opt2.InsecureHost {
		t.Error("SetInsecureHost option function not working")
	}
}
