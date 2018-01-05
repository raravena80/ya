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
	"bytes"
	"fmt"
	"github.com/raravena80/ya/common"
	"golang.org/x/crypto/ssh"
)

func executeCmd(opt common.Options, hostname string, config *ssh.ClientConfig) executeResult {

	port := fmt.Sprintf("%v", opt.Port)
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)

	if err != nil {
		return executeResult{
			result: hostname + ":\n",
			err:    err,
		}
	}

	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(opt.Cmd)

	return makeExecResult(hostname, stdoutBuf.String(), err)
}
