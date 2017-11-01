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

package shell

import (
	"bytes"
	"fmt"
	"github.com/raravena80/ya/common"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

type executeResult struct {
	result string
	err    error
}

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

func makeKeyring(key string, useAgent bool) ssh.AuthMethod {
	signers := []ssh.Signer{}

	if useAgent == true {
		aConn, _ := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
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
	return ssh.PublicKeys(signers...)
}

func executeCmd(opt common.Options, hostname string, config *ssh.ClientConfig) executeResult {

	conn, err := ssh.Dial("tcp", hostname+":"+opt.Port, config)

	if err != nil {
		return executeResult{result: "",
			err: err}
	}

	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(opt.Cmd)

	return executeResult{result: hostname + ":\n" + stdoutBuf.String(),
		err: err}
}

func Run(options ...func(*common.Options)) bool {
	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// in t seconds the message will come to timeout channel
	t, _ := strconv.ParseInt(opt.Timeout, 10, 64)
	timeout := time.After(time.Duration(t) * time.Second)
	results := make(chan executeResult, len(opt.Machines)+1)

	config := &ssh.ClientConfig{
		User: opt.User,
		Auth: []ssh.AuthMethod{
			makeKeyring(opt.Key, opt.UseAgent),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	for _, m := range opt.Machines {
		go func(hostname string) {
			results <- executeCmd(opt, hostname, config)
			// we’ll write results into the buffered channel of strings
		}(m)
	}

	retval := true

	for i := 0; i < len(opt.Machines); i++ {
		select {
		case res := <-results:
			if res.err == nil {
				fmt.Print(res.result)
			} else {
				fmt.Println(res.err)
				retval = false
			}
		case <-timeout:
			fmt.Println("Timed out!")
			retval = false
		}
	}
	return retval
}
