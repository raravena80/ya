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

func processDir(srcPath string, srcFileInfo os.FileInfo, procWriter, errPipe io.Writer) error {

	err := sendDir(srcPath, srcFileInfo, procWriter, errPipe)
	if err != nil {
		return err
	}
	dir, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			err = processDir(filepath.Join(srcPath, fi.Name()), fi, procWriter, errPipe)
			if err != nil {
				return err
			}
		} else {
			err = sendFile(filepath.Join(srcPath, fi.Name()), fi, procWriter, errPipe)
			if err != nil {
				return err
			}
		}
	}
	err = sendEndDir(procWriter, errPipe)
	return err
}

func sendEndDir(procWriter, errPipe io.Writer) error {
	header := fmt.Sprintf("E\n")
	_, err := procWriter.Write([]byte(header))
	return err
}

func sendDir(srcPath string, srcFileInfo os.FileInfo, procWriter, errPipe io.Writer) error {
	mode := uint32(srcFileInfo.Mode().Perm())
	header := fmt.Sprintf("D%04o 0 %s\n", mode, filepath.Base(srcPath))
	_, err := procWriter.Write([]byte(header))
	return err
}

func sendByte(w io.Writer, val byte) error {
	_, err := w.Write([]byte{val})
	return err
}

func sendFile(srcFile string, srcFileInfo os.FileInfo, procWriter, errPipe io.Writer) error {

	mode := uint32(srcFileInfo.Mode().Perm())
	fileReader, err := os.Open(srcFile)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not open source file "+srcFile, err.Error())
		return err
	}
	defer fileReader.Close()

	size := srcFileInfo.Size()
	header := fmt.Sprintf("C%04o %d %s\n", mode, size, filepath.Base(srcFile))

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

	srcFileInfo, err := os.Stat(opt.Src)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not stat source file "+opt.Src)
		return executeResult{result: hostname + ":\n", err: err}
	}

	if opt.IsRecursive {
		if srcFileInfo.IsDir() {
			err = processDir(opt.Src, srcFileInfo, procWriter, errPipe)
		} else {
			err = sendFile(opt.Src, srcFileInfo, procWriter, errPipe)
		}
	} else {
		if srcFileInfo.IsDir() {
			fmt.Fprintln(errPipe, "Not a regular file:", opt.Src, "specify recursive")
			err = fmt.Errorf("Not a regular file %v", opt.Src)
		} else {
			err = sendFile(opt.Src, srcFileInfo, procWriter, errPipe)
		}
	}

	return executeResult{result: hostname + ":\n",
		err: err}
}
