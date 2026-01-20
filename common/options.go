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

// Options holds the configuration for SSH/SCP operations.
// It contains connection details, authentication settings, and operation-specific parameters.
type Options struct {
	Machines     []string
	Port         int
	Timeout      int
	User         string
	Cmd          string
	Key          string
	Src          string
	Dst          string
	AgentSock    string
	Op           string
	UseAgent     bool
	IsRecursive  bool
	IsVerbose    bool
	KnownHosts   string
	InsecureHost bool
}

// SetUser Sets user for ssh session
func SetUser(u string) func(*Options) {
	return func(e *Options) {
		e.User = u
	}
}

// SetPort Sets port for ssh session
func SetPort(p int) func(*Options) {
	return func(e *Options) {
		e.Port = p
	}
}

// SetCmd Sets the command to be run on the ssh session
func SetCmd(c string) func(*Options) {
	return func(e *Options) {
		e.Cmd = c
	}
}

// SetMachines Sets the machines that we are going to run
// command or copy a file to
func SetMachines(m []string) func(*Options) {
	return func(e *Options) {
		e.Machines = m
	}
}

// SetKey Sets the key we are going to use to ssh connect
func SetKey(k string) func(*Options) {
	return func(e *Options) {
		e.Key = k
	}
}

// SetSource Sets the source file/dir for scp
func SetSource(s string) func(*Options) {
	return func(e *Options) {
		e.Src = s
	}
}

// SetDestination Sets the destination file/dir for scp
func SetDestination(d string) func(*Options) {
	return func(e *Options) {
		e.Dst = d
	}
}

// SetUseAgent Sets whether we want to use the ssh agent
func SetUseAgent(u bool) func(*Options) {
	return func(e *Options) {
		e.UseAgent = u
	}
}

// SetTimeout Sets the connection timeout
func SetTimeout(t int) func(*Options) {
	return func(e *Options) {
		e.Timeout = t
	}
}

// SetAgentSock Sets the ssh agent socket
func SetAgentSock(as string) func(*Options) {
	return func(e *Options) {
		e.AgentSock = as
	}
}

// SetOp Sets whether we want to run a command or scp
func SetOp(o string) func(*Options) {
	return func(e *Options) {
		e.Op = o
	}
}

// SetIsRecursive Sets whether we want to a recursive scp
func SetIsRecursive(r bool) func(*Options) {
	return func(e *Options) {
		e.IsRecursive = r
	}
}

// SetVerbose Sets high verbosity
func SetVerbose(v bool) func(*Options) {
	return func(e *Options) {
		e.IsVerbose = v
	}
}

// SetKnownHosts Sets the known_hosts file path for SSH host key verification
func SetKnownHosts(k string) func(*Options) {
	return func(e *Options) {
		e.KnownHosts = k
	}
}

// SetInsecureHost Sets whether to skip host key verification (not recommended for production)
func SetInsecureHost(i bool) func(*Options) {
	return func(e *Options) {
		e.InsecureHost = i
	}
}
