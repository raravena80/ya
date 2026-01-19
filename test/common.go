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

package test

import (
	"fmt"
	glssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"unsafe"
)

// MockSSHKey Mock SSH key for tests
type MockSSHKey struct {
	Keyname string
	Content []byte
	Privkey agent.AddedKey
	Pubkey  ssh.PublicKey
}

// SetupSSHAgent Setup and SSH agent for tests
func SetupSSHAgent(socketFile string) {
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
				fmt.Printf("Couldn't accept connection to agent tests %v\n", err)
			}
			go func(c io.ReadWriter) {
				err := agent.ServeAgent(a, c)
				if err != nil {
					fmt.Printf("Couldn't serve ssh agent for tests %v\n", err)
				}

			}(c)
		}

	}(done)
	<-done
}

// AddKeytoSSHAgent Adds a private key to the ssh agent
func AddKeytoSSHAgent(key agent.AddedKey, s string) {
	aConn, _ := net.Dial("unix", s)
	sshAgent := agent.NewClient(aConn)
	sshAgent.Add(key)
}

// RemoveKeyfromSSHAgent Removes a key from the ssh agent
func RemoveKeyfromSSHAgent(key ssh.PublicKey, s string) {
	aConn, _ := net.Dial("unix", s)
	sshAgent := agent.NewClient(aConn)
	sshAgent.Remove(key)
}

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// StartSSHServerForSSH Starts SSH server for ssh tests
func StartSSHServerForSSH(publicKeys map[string]ssh.PublicKey) {
	done := make(chan bool, 1)
	go func(done chan<- bool) {
		sshHandler := func(s glssh.Session) {
			authorizedKey := ssh.MarshalAuthorizedKey(s.PublicKey())
			io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
			s.Write(authorizedKey)
			io.WriteString(s, fmt.Sprintf("Command used %s:\n", s.Command()))
		}

		publicKeyOption := glssh.PublicKeyAuth(func(ctx glssh.Context, key glssh.PublicKey) bool {
			for _, pubk := range publicKeys {
				if glssh.KeysEqual(key, pubk) {
					return true
				}
			}
			return false // use glssh.KeysEqual() to compare against known keys
		})

		fmt.Println("starting ssh server for ssh tests on port 2222...")
		done <- true
		panic(glssh.ListenAndServe(":2222", sshHandler, publicKeyOption))
	}(done)
	<-done
}

// StartSSHServerForScp Starts SSH server for scp tests
func StartSSHServerForScp(publicKeys map[string]ssh.PublicKey) {
	done := make(chan bool, 1)
	go func(done chan<- bool) {
		scpHandler := func(s glssh.Session) {
			authorizedKey := ssh.MarshalAuthorizedKey(s.PublicKey())
			io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
			s.Write(authorizedKey)
			io.WriteString(s, fmt.Sprintf("Command used %s:\n", s.Command()))
			// Handle scp
			rp := regexp.MustCompile("scp")
			if rp.MatchString(s.Command()[0]) {
				cmd := exec.Command(s.Command()[0], s.Command()[1:]...)
				f, _ := cmd.StdinPipe()
				err := cmd.Start()
				if err != nil {
					panic(err)
				}
				go func() {
					io.Copy(f, s) // stdin
				}()
			}
		}

		publicKeyOption := glssh.PublicKeyAuth(func(ctx glssh.Context, key glssh.PublicKey) bool {
			for _, pubk := range publicKeys {
				if glssh.KeysEqual(key, pubk) {
					return true
				}
			}
			return false // use glssh.KeysEqual() to compare against known keys
		})

		fmt.Println("starting ssh server for scp tests on port 2224...")
		done <- true
		panic(glssh.ListenAndServe(":2224", scpHandler, publicKeyOption))
	}(done)
	<-done
}
