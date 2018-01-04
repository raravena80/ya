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
	"golang.org/x/crypto/ssh"
	"time"
)

type executeResult struct {
	result string
	err    error
}

type execFuncType func(common.Options, string, *ssh.ClientConfig) executeResult

// SSHSession Create an SSH Session
func SSHSession(options ...func(*common.Options)) bool {
	var execFunc execFuncType

	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// in opt.Timeout seconds the message will come to timeout channel
	done := make(chan bool, len(opt.Machines))

	sshAuth := []ssh.AuthMethod{
		ssh.PublicKeys(common.MakeKeyring(
			opt.Key,
			opt.AgentSock,
			opt.UseAgent)...),
	}
	config := &ssh.ClientConfig{
		User:            opt.User,
		Auth:            sshAuth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(opt.Timeout) * time.Second,
	}

	for _, m := range opt.Machines {
		// we’ll write results into the buffered channel of strings
		switch opt.Op {
		case "ssh":
			execFunc = executeCmd
		case "scp":
			execFunc = executeCopy
		}
		go func(hostname string, execFunc execFuncType) {
			res := execFunc(opt, hostname, config)
			if res.err == nil {
				fmt.Print(res.result)
				done <- true
			} else {
				fmt.Println(res.result, "\n", res.err)
				done <- false
			}
		}(m, execFunc)
	}

	retval := true
	for i := 0; i < len(opt.Machines); i++ {
		if !<-done {
			retval = false
		}
	}
	return retval
}
