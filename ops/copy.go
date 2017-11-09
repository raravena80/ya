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
	"io"
	"os"
	"path/filepath"
)

func sendByte(w io.Writer, val byte) error {
	_, err := w.Write([]byte{val})
	return err
}

func sendFile(opt common.Options, procWriter, errPipe io.Writer) error {

	srcFileInfo, err := os.Stat(opt.Src)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not stat source file "+opt.Src)
		return err
	}
	mode := uint32(srcFileInfo.Mode().Perm())

	fileReader, err := os.Open(opt.Src)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not open source file "+opt.Src, err.Error())
		return err
	}
	defer fileReader.Close()

	size := srcFileInfo.Size()
	header := fmt.Sprintf("C%04o %d %s\n", mode, size, filepath.Base(opt.Dst))

	_, err = procWriter.Write([]byte(header))
	if err != nil {
		fmt.Fprintln(errPipe, "Could not write scp header", err.Error())
		return err
	}

	_, err = io.Copy(procWriter, fileReader)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not send file", err.Error())
		return err
	}
	// terminate with null byte
	err = sendByte(procWriter, 0)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not send last byte", err.Error())
		return err
	}

	err = fileReader.Close()
	if err != nil {
		fmt.Fprintln(errPipe, "Could not close source file "+opt.Src, err.Error())
	}
	return err
}

func executeCopy(opt common.Options, hostname string, config *ssh.ClientConfig) executeResult {

	port := fmt.Sprintf("%v", opt.Port)
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)

	if err != nil {
		return executeResult{
			result: "Connection error",
			err:    err,
		}
	}
	session, _ := conn.NewSession()
	defer session.Close()

	errPipe := os.Stderr
	procWriter, err := session.StdinPipe()

	if err != nil {
		fmt.Fprintln(errPipe, err.Error())
	}
	defer procWriter.Close()

	scpCmd := fmt.Sprintf("/usr/bin/scp -qrt %s/.", filepath.Dir(opt.Dst))
	err = session.Start(scpCmd)
	if err != nil {
		fmt.Fprintln(errPipe, err.Error())
	}
	err = sendFile(opt, procWriter, errPipe)
	return executeResult{result: hostname + ":\n",
		err: err}
}
