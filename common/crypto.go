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
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
)

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

// MakeKeyring Makes an ssh key ring for authentication
func MakeKeyring(key, agentSock string, useAgent bool) []ssh.Signer {
	signers := []ssh.Signer{}

	if useAgent == true {
		aConn, _ := net.Dial("unix", agentSock)
		sshAgent := agent.NewClient(aConn)
		aSigners, _ := sshAgent.Signers()
		for _, signer := range aSigners {
			signers = append(signers, signer)
		}
	}

	keys := []string{key}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}
	return signers
}
