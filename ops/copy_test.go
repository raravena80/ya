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
	"bytes"
	"fmt"
	"github.com/raravena80/ya/common"
	"github.com/raravena80/ya/test"
	"golang.org/x/crypto/ssh/testdata"
	"os"
	"testing"
	"time"
)

func init() {
	// Create files to test scp locally
	os.WriteFile("/tmp/removethis1", []byte("Sample file 1  uu"), 0644)
	os.WriteFile("/tmp/removethisnoperm", []byte("Sample file 2 uu"), 0000)
	os.Mkdir("/tmp/removethisdir", 0777)
	os.Mkdir("/tmp/removethisdir/removethisotherdir", 0777)
	os.WriteFile("/tmp/removethisdir/removefile", []byte("Sample file in dir"), 0644)
	os.WriteFile("/tmp/removethisdir/removethisotherdir/file1", []byte("another_file"), 0644)
	os.Mkdir("/tmp/anotherremovethisdir", 0777)
	os.Mkdir("/tmp/anotherremovethisdir/second", 0777)
	os.WriteFile("/tmp/anotherremovethisdir/second/file1", []byte("another_file"), 000)
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		machines    []string
		port        int
		timeout     int
		user        string
		cmd         string
		key         test.MockSSHKey
		op          string
		src         string
		dst         string
		useagent    bool
		isrecursive bool
		verbose     bool
		expected    bool
	}{
		{name: "Basic with valid rsa key wrong hostname",
			machines: []string{"bogushost"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey12",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: false,
		},
		{name: "Basic with valid rsa key wrong port",
			machines: []string{"localhost"},
			port:     2223,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey13",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: false,
		},
		{name: "Basic with valid rsa key Google endpoint",
			machines: []string{"www.google.com"},
			port:     22,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey14",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  1,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: false,
		},
		{name: "Basic with valid rsa key scp",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey15",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethis1",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: true,
		},
		{name: "Basic with valid rsa key scp wrong file",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey16",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/doesntexist",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: false,
		},
		{name: "Basic with valid rsa key scp file with no permissions",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey16",
				Content: testdata.PEMBytes["rsa"],
			},
			op:       "scp",
			useagent: false,
			timeout:  5,
			src:      "/tmp/removethisnoperm",
			dst:      "/tmp/removethis2",
			verbose:  true,
			expected: false,
		},
		{name: "Basic with valid rsa key scp dir recursive",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey17",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/removethisdir",
			dst:         "/tmp/removethisdir2",
			expected:    true,
			verbose:     true,
			isrecursive: true,
		},
		{name: "Basic with valid rsa key scp dir non recursive",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey18",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/removethisdir",
			dst:         "/tmp/removethisdir3",
			expected:    false,
			verbose:     true,
			isrecursive: false,
		},
		{name: "Basic with valid rsa key scp file recursive",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey19",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/removethis1",
			dst:         "/tmp/wedontcare",
			expected:    true,
			verbose:     true,
			isrecursive: true,
		},
		{name: "Basic with valid rsa key scp dir recursive file no permissions",
			machines: []string{"127.0.0.1"},
			port:     2224,
			cmd:      "ls",
			user:     "testuser",
			key: test.MockSSHKey{
				Keyname: "/tmp/mockkey20",
				Content: testdata.PEMBytes["rsa"],
			},
			op:          "scp",
			useagent:    false,
			timeout:     5,
			src:         "/tmp/anotherremovethisdir",
			dst:         "/tmp/anotherremovethisdir2",
			expected:    false,
			verbose:     true,
			isrecursive: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write Content of the key to the Keyname file
			if tt.key.Keyname != "" {
				os.WriteFile(tt.key.Keyname, tt.key.Content, 0644)
			}
			returned := SSHSession(common.SetMachines(tt.machines),
				common.SetUser(tt.user),
				common.SetPort(tt.port),
				common.SetCmd(tt.cmd),
				common.SetKey(tt.key.Keyname),
				common.SetUseAgent(tt.useagent),
				common.SetTimeout(tt.timeout),
				common.SetSource(tt.src),
				common.SetDestination(tt.dst),
				common.SetOp(tt.op),
				common.SetIsRecursive(tt.isrecursive),
				common.SetVerbose(tt.verbose))

			if !(returned == tt.expected) {
				t.Errorf("Value received: %v expected %v", returned, tt.expected)
			}
			if tt.key.Keyname != "" {
				os.Remove(tt.key.Keyname)
			}
		})
	}
}

func TestSendEndDir(t *testing.T) {
	var buf bytes.Buffer
	err := sendEndDir(&buf, &buf)
	if err != nil {
		t.Errorf("sendEndDir() returned error: %v", err)
	}
	result := buf.String()
	if result != "E\n" {
		t.Errorf("sendEndDir() = %q, want %q", result, "E\n")
	}
}

func TestSendDir(t *testing.T) {
	tests := []struct {
		name         string
		srcPath      string
		mode         os.FileMode
		expectedSub  string
	}{
		{name: "Send directory with 755 permissions",
			srcPath: "/tmp/testdir",
			mode: 0755,
			expectedSub: "0755"},
		{name: "Send directory with 644 permissions",
			srcPath: "/path/to/mydir",
			mode: 0644,
			expectedSub: "0644"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			fi, err := os.Stat(tt.srcPath)
			if err != nil {
				// Create temp file info for testing
				fi, _ = os.Lstat(".")
			}
			// Use test mode instead of actual file mode
			fi = &mockFileInfo{name: tt.srcPath, mode: tt.mode}

			err = sendDir(tt.srcPath, fi, &buf, &buf)
			if err != nil {
				t.Errorf("sendDir() returned error: %v", err)
			}
			result := buf.String()
			if !contains(result, tt.expectedSub) {
				t.Errorf("sendDir() = %q, want to contain %q", result, tt.expectedSub)
			}
		})
	}
}

func TestSendByte(t *testing.T) {
	tests := []struct {
		name     string
		val      byte
		expected []byte
	}{
		{name: "Send null byte", val: 0, expected: []byte{0}},
		{name: "Send byte 65", val: 65, expected: []byte{65}},
		{name: "Send byte 255", val: 255, expected: []byte{255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := sendByte(&buf, tt.val)
			if err != nil {
				t.Errorf("sendByte() returned error: %v", err)
			}
			result := buf.Bytes()
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("sendByte() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		message     string
		verbose     bool
		expectError bool
	}{
		{name: "Error with verbose true",
			err: fmt.Errorf("test error"),
			message: "Test message: ",
			verbose: true,
			expectError: true},
		{name: "Error with verbose false",
			err: fmt.Errorf("test error"),
			message: "Test message: ",
			verbose: false,
			expectError: true},
		{name: "Nil error with verbose true",
			err: nil,
			message: "Test message",
			verbose: true,
			expectError: false},
		{name: "Nil error with verbose false",
			err: nil,
			message: "Test message",
			verbose: false,
			expectError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			result := processError(tt.err, tt.message, &buf, tt.verbose)
			if tt.expectError && result == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && result != nil {
				t.Errorf("Expected nil error, got %v", result)
			}
			if tt.verbose && tt.err != nil {
				output := buf.String()
				if !contains(output, tt.message) {
					t.Errorf("Expected output to contain %q, got %q", tt.message, output)
				}
			}
		})
	}
}

func TestDefaultSCPPath(t *testing.T) {
	if DefaultSCPPath != "/usr/bin/scp" {
		t.Errorf("DefaultSCPPath = %s, want /usr/bin/scp", DefaultSCPPath)
	}
}

// Helper functions
type mockFileInfo struct {
	name string
	mode os.FileMode
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.mode.IsDir() }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestTearCopy(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "Teardown Copy test",
			id: "copyTestTdown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == "copyTestTdown" {
				os.Remove("/tmp/removethis1")
				os.Remove("/tmp/removethisnoperm")
				os.RemoveAll("/tmp/removethisdir")
				os.RemoveAll("/tmp/removethisdir2")
				os.RemoveAll("/tmp/anotherremovethisdir/")
			}

		})

	}
}
