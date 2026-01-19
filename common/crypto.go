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
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
)

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file %s: %w", keyname, err)
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %s: %w", keyname, err)
	}

	signer, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key %s: %w", keyname, err)
	}
	return signer, nil
}

// MakeKeyring Makes an ssh key ring for authentication
func MakeKeyring(key, agentSock string, useAgent bool) []ssh.Signer {
	signers := []ssh.Signer{}

	if useAgent {
		aConn, err := net.Dial("unix", agentSock)
		if err == nil {
			defer aConn.Close()
			sshAgent := agent.NewClient(aConn)
			aSigners, err := sshAgent.Signers()
			if err == nil {
				for _, signer := range aSigners {
					signers = append(signers, signer)
				}
			}
		}
		// Continue with key-based auth even if agent fails
	}

	keys := []string{key}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
		// Log error but don't fail - user may have provided an invalid key
		if err != nil && keyname != "" {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}

	if len(signers) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: No valid SSH authentication methods available")
	}

	return signers
}
