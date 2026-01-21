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

func TestSetConnectTimeout(t *testing.T) {
	tests := []struct {
		name           string
		connectTimeout int
		expected       *int
	}{
		{name: "Set connect timeout to 30", connectTimeout: 30},
		{name: "Set connect timeout to 60", connectTimeout: 60},
		{name: "Set connect timeout to 0", connectTimeout: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetConnectTimeout(tt.connectTimeout)(&opt)
			if opt.ConnectTimeout == nil {
				if tt.connectTimeout != 0 {
					t.Error("ConnectTimeout should not be nil")
				}
			} else if *opt.ConnectTimeout != tt.connectTimeout {
				t.Errorf("SetConnectTimeout() = %v, want %v", *opt.ConnectTimeout, tt.connectTimeout)
			}
		})
	}
}

func TestSetCommandTimeout(t *testing.T) {
	tests := []struct {
		name            string
		commandTimeout  int
		expected        *int
	}{
		{name: "Set command timeout to 30", commandTimeout: 30},
		{name: "Set command timeout to 120", commandTimeout: 120},
		{name: "Set command timeout to 0", commandTimeout: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetCommandTimeout(tt.commandTimeout)(&opt)
			if opt.CommandTimeout == nil {
				if tt.commandTimeout != 0 {
					t.Error("CommandTimeout should not be nil")
				}
			} else if *opt.CommandTimeout != tt.commandTimeout {
				t.Errorf("SetCommandTimeout() = %v, want %v", *opt.CommandTimeout, tt.commandTimeout)
			}
		})
	}
}

func TestSetOutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{name: "Set format to text", format: "text", expected: "text"},
		{name: "Set format to json", format: "json", expected: "json"},
		{name: "Set format to yaml", format: "yaml", expected: "yaml"},
		{name: "Set format to table", format: "table", expected: "table"},
		{name: "Set format to empty", format: "", expected: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetOutputFormat(tt.format)(&opt)
			if opt.OutputFormat != tt.expected {
				t.Errorf("SetOutputFormat() = %v, want %v", opt.OutputFormat, tt.expected)
			}
		})
	}
}

func TestSetDryRun(t *testing.T) {
	tests := []struct {
		name     string
		dryRun   bool
		expected bool
	}{
		{name: "Set dry-run to true", dryRun: true, expected: true},
		{name: "Set dry-run to false", dryRun: false, expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetDryRun(tt.dryRun)(&opt)
			if opt.DryRun != tt.expected {
				t.Errorf("SetDryRun() = %v, want %v", opt.DryRun, tt.expected)
			}
		})
	}
}

func TestSetHostPatterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		expected []string
	}{
		{name: "Set single pattern", patterns: []string{"*prod*"}, expected: []string{"*prod*"}},
		{name: "Set multiple patterns", patterns: []string{"*prod*", "*staging*"}, expected: []string{"*prod*", "*staging*"}},
		{name: "Set empty patterns", patterns: []string{}, expected: []string{}},
		{name: "Set nil patterns", patterns: nil, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetHostPatterns(tt.patterns)(&opt)
			if len(opt.HostPatterns) != len(tt.expected) {
				t.Errorf("SetHostPatterns() length = %v, want %v", len(opt.HostPatterns), len(tt.expected))
			}
			for i, p := range opt.HostPatterns {
				if p != tt.expected[i] {
					t.Errorf("SetHostPatterns()[%d] = %v, want %v", i, p, tt.expected[i])
				}
			}
		})
	}
}

func TestSetHostExcludes(t *testing.T) {
	tests := []struct {
		name     string
		excludes []string
		expected []string
	}{
		{name: "Set single exclude", excludes: []string{"*backup*"}, expected: []string{"*backup*"}},
		{name: "Set multiple excludes", excludes: []string{"*backup*", "*test*"}, expected: []string{"*backup*", "*test*"}},
		{name: "Set empty excludes", excludes: []string{}, expected: []string{}},
		{name: "Set nil excludes", excludes: nil, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetHostExcludes(tt.excludes)(&opt)
			if len(opt.HostExcludes) != len(tt.expected) {
				t.Errorf("SetHostExcludes() length = %v, want %v", len(opt.HostExcludes), len(tt.expected))
			}
			for i, e := range opt.HostExcludes {
				if e != tt.expected[i] {
					t.Errorf("SetHostExcludes()[%d] = %v, want %v", i, e, tt.expected[i])
				}
			}
		})
	}
}

func TestSetShowProgress(t *testing.T) {
	tests := []struct {
		name     string
		progress bool
		expected bool
	}{
		{name: "Set show progress to true", progress: true, expected: true},
		{name: "Set show progress to false", progress: false, expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{}
			SetShowProgress(tt.progress)(&opt)
			if opt.ShowProgress != tt.expected {
				t.Errorf("SetShowProgress() = %v, want %v", opt.ShowProgress, tt.expected)
			}
		})
	}
}
