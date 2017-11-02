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

type Options struct {
	Machines  []string
	Port      string
	User      string
	Cmd       string
	Key       string
	Src       string
	Dst       string
	Timeout   string
	AgentSock string
	UseAgent  bool
}

func SetUser(u string) func(*Options) {
	return func(e *Options) {
		e.User = u
	}
}

func SetPort(p string) func(*Options) {
	return func(e *Options) {
		e.Port = p
	}
}

func SetCmd(c string) func(*Options) {
	return func(e *Options) {
		e.Cmd = c
	}
}

func SetMachines(m []string) func(*Options) {
	return func(e *Options) {
		e.Machines = m
	}
}

func SetKey(k string) func(*Options) {
	return func(e *Options) {
		e.Key = k
	}
}

func SetSource(s string) func(*Options) {
	return func(e *Options) {
		e.Src = s
	}
}

func SetDestination(d string) func(*Options) {
	return func(e *Options) {
		e.Dst = d
	}
}

func SetUseAgent(u bool) func(*Options) {
	return func(e *Options) {
		e.UseAgent = u
	}
}

func SetTimeout(t string) func(*Options) {
	return func(e *Options) {
		e.Timeout = t
	}
}

func SetAgentSock(as string) func(*Options) {
	return func(e *Options) {
		e.AgentSock = as
	}
}
