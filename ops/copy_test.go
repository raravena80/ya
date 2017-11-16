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
package ops

import (
	"github.com/raravena80/ya/common"
	"github.com/raravena80/ya/test"
	"golang.org/x/crypto/ssh/testdata"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	// Create files to test scp locally
	ioutil.WriteFile("/tmp/removethis1", []byte("Sample file 1  uu"), 0644)
	ioutil.WriteFile("/tmp/removethisnoperm", []byte("Sample file 2 uu"), 0000)
	os.Mkdir("/tmp/removethisdir", 0777)
	ioutil.WriteFile("/tmp/removethisdir/removefile", []byte("Sample file in dir"), 0644)
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		machines    []string
		port        int
		timeout     int
		user        string
		cmd         string
		key         test.MockSSHKey
		op          string
		src         string
		dst         string
		useagent    bool
		isrecursive bool
		verbose     bool
		expected    bool
	}{
		{name: "Basic with valid rsa key wrong hostname",
			machines: []string{"bogushost"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey12",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			expected: false,
		},
		{name: "Basic with valid rsa key wrong port",
			machines: []string{"localhost"},
			port:     2223,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey13",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			expected: false,
		},
		{name: "Basic with valid rsa key Google endpoint",
			machines: []string{"www.google.com"},
			port:     22,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey14",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  1,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			expected: false,
		},
		{name: "Basic with valid rsa key scp",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey15",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			expected: true,
		},
		{name: "Basic with valid rsa key scp wrong file",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey16",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/doesntexist",
			dst:      "/tmp/removethis2",
			expected: false,
		},
		{name: "Basic with valid rsa key scp file with no permissions",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey16",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethisnoperm",
			dst:      "/tmp/removethis2",
			expected: false,
		},
		{name: "Basic with valid rsa key scp dir recursive",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey17",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/removethisdir",
			dst:         "/tmp/removethisdir",
			expected:    true,
			isrecursive: true,
		},
		{name: "Basic with valid rsa key scp dir non recursive",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey18",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/removethisdir",
			dst:         "/tmp/removethisdir2",
			expected:    true,
			isrecursive: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write Content of the key to the Keyname file
			if tt.key.Keyname != "" {
				ioutil.WriteFile(tt.key.Keyname, tt.key.Content, 0644)
			}
			returned := SSHSession(common.SetMachines(tt.machines),
				common.SetUser(tt.user),
				common.SetPort(tt.port),
				common.SetCmd(tt.cmd),
				common.SetKey(tt.key.Keyname),
				common.SetUseAgent(tt.useagent),
				common.SetTimeout(tt.timeout),
				common.SetSource(tt.src),
				common.SetDestination(tt.dst),
				common.SetOp(tt.op),
				common.SetIsRecursive(tt.isrecursive),
				common.SetVerbose(tt.verbose))

			if !(returned == tt.expected) {
				t.Errorf("Value received: %v expected %v", returned, tt.expected)
			}
			if tt.key.Keyname != "" {
				os.Remove(tt.key.Keyname)
			}
		})
	}
}

func TestTearCopy(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "Teardown Copy test",
			id: "copyTestTdown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == "copyTestTdown" {
				os.Remove("/tmp/removethis1")
				os.Remove("/tmp/removethisnoperm")
				os.Remove("/tmp/removethisdir/removefile")
				os.Remove("/tmp/removethisdir")
			}

		})

	}
}
