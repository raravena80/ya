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
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	// maxKeySize is the maximum allowed size for an SSH key file (1MB)
	maxKeySize = 1024 * 1024
)

// makeSigner creates an SSH signer from a private key file.
// It validates the key file permissions and size before parsing.
func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file %s: %w", keyname, err)
	}
	defer fp.Close()

	// Check file permissions for security
	info, err := fp.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat key file %s: %w", keyname, err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 && perm != 0400 {
		fmt.Fprintf(os.Stderr, "Warning: SSH key file has insecure permissions: %03o (expected 0600 or 0400)\n", perm)
	}

	// Limit the read size to prevent potential issues with very large files
	limitedReader := io.LimitReader(fp, maxKeySize)
	buf, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %s: %w", keyname, err)
	}

	signer, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key %s: %w", keyname, err)
	}
	return signer, nil
}

// MakeKeyring creates an SSH keyring for authentication.
// It attempts to use the SSH agent if useAgent is true, and also tries to load
// the key from the specified key file. Returns a slice of SSH signers.
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
