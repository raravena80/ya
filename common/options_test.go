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

package common

import (
	"testing"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name      string
		machines  []string
		port      int
		user      string
		cmd       string
		key       string
		src       string
		dst       string
		timeout   string
		agentSock string
		useAgent  bool
	}{
		{name: "Test all options",
			machines:  []string{"one", "two", "three"},
			port:      22,
			user:      "bogus",
			cmd:       "runit",
			key:       "mykey",
			src:       "src",
			dst:       "dst",
			timeout:   "20",
			agentSock: "socket",
			useAgent:  false,
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
				SetUseAgent(tt.useAgent)}
			opt := Options{}
			for _, option := range options {
				option(&opt)
			}
		})
	}
}
