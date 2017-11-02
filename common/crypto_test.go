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
	"fmt"
	"github.com/raravena80/ya/test"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/testdata"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

var (
	testPrivateKeys map[string]interface{}
	testSigners     map[string]ssh.Signer
	testPublicKeys  map[string]ssh.PublicKey
	sshAgentSocket  string
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

	randomStr := fmt.Sprintf("%v", rand.Intn(5000))
	sshAgentSocket = "/tmp/gosocket" + randomStr + ".sock"
	test.SetupSshAgent(sshAgentSocket)
}

func TestMakeSigner(t *testing.T) {
	tests := []struct {
		name     string
		key      test.MockSSHKey
		expected ssh.Signer
	}{
		{name: "Basic key signer with valid rsa key",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			expected: testSigners["rsa"],
		},
		{name: "Basic key signer with valid dsa key",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["dsa"],
			},
			expected: testSigners["dsa"],
		},
		{name: "Basic key signer with valid ecdsa key",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["ecdsa"],
			},
			expected: testSigners["ecdsa"],
		},
		{name: "Basic key signer with valid user key",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["user"],
			},
			expected: testSigners["user"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write Content of the key to the Keyname file
			ioutil.WriteFile(tt.key.Keyname, tt.key.Content, 0644)
			returned, _ := makeSigner(tt.key.Keyname)
			if !reflect.DeepEqual(returned, tt.expected) {
				t.Errorf("Value received: %v expected %v", returned, tt.expected)
			}
			os.Remove(tt.key.Keyname)
		})
	}
}

func TestMakeKeyring(t *testing.T) {
	tests := []struct {
		name     string
		useagent bool
		key      test.MockSSHKey
		expected []byte
	}{
		{name: "Basic key ring with valid rsa key",
			useagent: false,
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["rsa"],
			},
			expected: testPublicKeys["rsa"].Marshal(),
		},
		{name: "Basic key ring with valid dsa key",
			useagent: false,
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["dsa"],
			},
			expected: testPublicKeys["dsa"].Marshal(),
		},
		{name: "Basic key ring with valid ecdsa key",
			useagent: false,
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["ecdsa"],
			},
			expected: testPublicKeys["ecdsa"].Marshal(),
		},
		{name: "Basic key ring with valid user key",
			useagent: false,
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey",
				Content: testdata.PEMBytes["user"],
			},
			expected: testPublicKeys["user"].Marshal(),
		},
		{name: "Basic key ring agent with valid rsa key",
			useagent: true,
			key: test.MockSSHKey{
				Keyname: "",
				Content: testdata.PEMBytes["rsa"],
				Privkey: agent.AddedKey{PrivateKey: testPrivateKeys["rsa"]},
				Pubkey:  testPublicKeys["rsa"],
			},
			expected: testPublicKeys["rsa"].Marshal(),
		},
		{name: "Basic key ring agent with valid dsa key",
			useagent: true,
			key: test.MockSSHKey{
				Keyname: "",
				Content: testdata.PEMBytes["dsa"],
				Privkey: agent.AddedKey{PrivateKey: testPrivateKeys["dsa"]},
				Pubkey:  testPublicKeys["dsa"],
			},
			expected: testPublicKeys["dsa"].Marshal(),
		},
		{name: "Basic key ring agent with valid ecdsa key",
			useagent: true,
			key: test.MockSSHKey{
				Keyname: "",
				Content: testdata.PEMBytes["ecdsa"],
				Privkey: agent.AddedKey{PrivateKey: testPrivateKeys["ecdsa"]},
				Pubkey:  testPublicKeys["ecdsa"],
			},
			expected: testPublicKeys["ecdsa"].Marshal(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useagent == true {
				test.AddKeytoSSHAgent(tt.key.Privkey, sshAgentSocket)
			}
			// Write Content of the key to the Keyname file
			if tt.key.Keyname != "" {
				ioutil.WriteFile(tt.key.Keyname, tt.key.Content, 0644)
			}
			signers := MakeKeyring(tt.key.Keyname, sshAgentSocket, tt.useagent)
			returned := signers[0].PublicKey().Marshal()
			// DeepEqual always returns false for functions unless nil
			// hence converting to string to compare
			if !reflect.DeepEqual(returned, tt.expected) {
				t.Errorf("Value received: %v expected %v", returned, tt.expected)
			}
			if tt.useagent == true {
				test.RemoveKeyfromSSHAgent(tt.key.Pubkey, sshAgentSocket)
			}
			if tt.key.Keyname != "" {
				os.Remove(tt.key.Keyname)
			}
		})
	}
}

func TestTearDown(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "Teardown SSH Agent",
			id: "sshAgentTdown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == "sshAgentTdown" {
				os.Remove(sshAgentSocket)
			}

		})

	}
}
