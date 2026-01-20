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

package ops

import (
	"fmt"
	"os"
	"time"

	"github.com/raravena80/ya/common"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
)

type executeResult struct {
	result string
	err    error
}

type execFuncType func(common.Options, string, *ssh.ClientConfig) executeResult

// Makes a common execResult
func makeExecResult(hostname, output string, err error) executeResult {
	return executeResult{
		result: hostname + ":\n" + output,
		err:    err,
	}
}

// getHostKeyCallback returns an appropriate HostKeyCallback based on options.
// It uses known_hosts file for proper host key verification, with fallback to
// insecure mode when explicitly requested or when known_hosts is not available.
func getHostKeyCallback(opt common.Options) ssh.HostKeyCallback {
	// If insecure mode is explicitly requested, use the insecure callback with a warning
	if opt.InsecureHost {
		fmt.Fprintln(os.Stderr, "Warning: Using insecure host key verification. This is not recommended for production.")
		return ssh.InsecureIgnoreHostKey()
	}

	// Try to use the user's known_hosts file for proper host key verification
	knownHostsPath := os.ExpandEnv("~/.ssh/known_hosts")
	callback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		// If known_hosts file doesn't exist or is unreadable, fall back to insecure mode with warning
		fmt.Fprintln(os.Stderr, "Warning: Could not read known_hosts file:", err)
		fmt.Fprintln(os.Stderr, "Warning: Falling back to insecure host key verification.")
		fmt.Fprintln(os.Stderr, "To fix this, ensure ~/.ssh/known_hosts exists or use --insecure-host to suppress this warning.")
		return ssh.InsecureIgnoreHostKey()
	}

	// Wrap the knownhosts callback to ensure type compatibility
	return ssh.HostKeyCallback(callback)
}

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
		HostKeyCallback: getHostKeyCallback(opt),
		Timeout:         time.Duration(opt.Timeout) * time.Second,
	}

	for _, m := range opt.Machines {
		// we'll write results into the buffered channel of strings
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
