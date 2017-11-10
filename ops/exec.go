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

// SSHSession Create an SSH Session
func SSHSession(options ...func(*common.Options)) bool {
	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// in opt.Timeout seconds the message will come to timeout channel
	timeout := time.After(time.Duration(opt.Timeout) * time.Second)
	results := make(chan executeResult, len(opt.Machines)+1)

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
	}

	for _, m := range opt.Machines {
		go func(hostname string) {
			// we’ll write results into the buffered channel of strings
			switch opt.Op {
			case "ssh":
				results <- executeCmd(opt, hostname, config)
			case "scp":
				results <- executeCopy(opt, hostname, config)
			}
		}(m)
	}

	retval := true

	for i := 0; i < len(opt.Machines); i++ {
		select {
		case res := <-results:
			if res.err == nil {
				fmt.Print(res.result)
			} else {
				fmt.Println(res.result, "\n", res.err)
				retval = false
			}
		case <-timeout:
			fmt.Println("Timed out!")
			retval = false
		}
	}
	return retval
}
