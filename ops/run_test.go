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
	"fmt"
	"github.com/raravena80/ya/common"
	"github.com/raravena80/ya/test"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testPrivateKeys map[string]interface{}
	testSigners     map[string]ssh.Signer
	testPublicKeys  map[string]ssh.PublicKey
)

func init() {
	var err error

	n := len(testdata.PEMBytes)
	testSigners = make(map[string]ssh.Signer, n)
	testPrivateKeys = make(map[string]interface{}, n)
	testPublicKeys = make(map[string]ssh.PublicKey, n)

	for t, k := range testdata.PEMBytes {
		testPrivateKeys[t], err = ssh.ParseRawPrivateKey(k)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse test key %s: %v", t, err))
		}
		testSigners[t], err = ssh.NewSignerFromKey(testPrivateKeys[t])
		if err != nil {
			panic(fmt.Sprintf("Unable to create signer for test key %s: %v", t, err))
		}
		testPublicKeys[t] = testSigners[t].PublicKey()
	}

	test.StartSshServer(testPublicKeys)
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		machines []string
		port     int
		timeout  int
		user     string
		cmd      string
		key      test.MockSshKey
		useagent bool
		expected bool
	}{
		{name: "Basic with valid rsa key",
			machines: []string{"localhost"},
			port:     2222,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSshKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			useagent: false,
			timeout:  5,
			expected: true,
		},
		{name: "Basic with valid rsa key wrong hostname",
			machines: []string{"bogushost"},
			port:     2222,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSshKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			useagent: false,
			timeout:  5,
			expected: false,
		},
		{name: "Basic with valid rsa key wrong port",
			machines: []string{"localhost"},
			port:     2223,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSshKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			useagent: false,
			timeout:  5,
			expected: false,
		},
		{name: "Basic with valid rsa key Google endpoint",
			machines: []string{"www.google.com"},
			port:     22,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSshKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			useagent: false,
			timeout:  1,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write Content of the key to the Keyname file
			if tt.key.Keyname != "" {
				ioutil.WriteFile(tt.key.Keyname, tt.key.Content, 0644)
			}
			returned := Run(common.SetMachines(tt.machines),
				common.SetUser(tt.user),
				common.SetPort(tt.port),
				common.SetCmd(tt.cmd),
				common.SetKey(tt.key.Keyname),
				common.SetUseAgent(tt.useagent),
				common.SetTimeout(tt.timeout))

			if !(returned == tt.expected) {
				t.Errorf("Value received: %v expected %v", returned, tt.expected)
			}
			if tt.key.Keyname != "" {
				os.Remove(tt.key.Keyname)
			}
		})
	}
}
