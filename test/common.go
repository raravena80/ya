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

package test

import (
	"fmt"
	glssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
)

type MockSshKey struct {
	Keyname string
	Content []byte
	Privkey agent.AddedKey
	Pubkey  ssh.PublicKey
}

func SetupSshAgent(socketFile string) {
	done := make(chan string, 1)
	a := agent.NewKeyring()
	ln, err := net.Listen("unix", socketFile)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create socket for tests %v", err))
	}

	go func(done chan<- string) {
		// Need to wait until the socket is setup
		firstTime := true
		for {
			if firstTime == true {
				done <- socketFile
				firstTime = false
			}
			c, err := ln.Accept()
			defer c.Close()
			if err != nil {
				panic(fmt.Sprintf("Couldn't accept connection to agent tests %v", err))
			}
			go func(c io.ReadWriter) {
				err := agent.ServeAgent(a, c)
				if err != nil {
					panic(fmt.Sprintf("Couldn't serve ssh agent for tests %v", err))
				}

			}(c)
		}

	}(done)
	<-done
}

func AddKeytoSshAgent(key agent.AddedKey, s string) {
	aConn, _ := net.Dial("unix", s)
	sshAgent := agent.NewClient(aConn)
	sshAgent.Add(key)
}

func RemoveKeyfromSshAgent(key ssh.PublicKey, s string) {
	aConn, _ := net.Dial("unix", s)
	sshAgent := agent.NewClient(aConn)
	sshAgent.Remove(key)
}

func StartSshServer(publicKeys map[string]ssh.PublicKey) {
	done := make(chan bool, 1)
	go func(done chan<- bool) {
		glssh.Handle(func(s glssh.Session) {
			authorizedKey := ssh.MarshalAuthorizedKey(s.PublicKey())
			io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
			s.Write(authorizedKey)
		})

		publicKeyOption := glssh.PublicKeyAuth(func(ctx glssh.Context, key glssh.PublicKey) bool {
			for _, pubk := range publicKeys {
				if glssh.KeysEqual(key, pubk) {
					return true
				}
			}
			return false // use glssh.KeysEqual() to compare against known keys
		})

		fmt.Println("starting ssh server on port 2222...")
		done <- true
		panic(glssh.ListenAndServe(":2222", nil, publicKeyOption))
	}(done)
	<-done
}