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
	"github.com/raravena80/ya/common"
	"testing"
)

func TestMakeExecResult(t *testing.T) {
	result := makeExecResult("testhost", "output", nil)

	if result.result != "testhost:\noutput" {
		t.Errorf("Expected 'testhost:\\noutput', got '%s'", result.result)
	}
	if result.err != nil {
		t.Errorf("Expected nil error, got %v", result.err)
	}
}

func TestMakeExecResultWithError(t *testing.T) {
	testErr := fmt.Errorf("connection closed")
	result := makeExecResult("testhost", "output", testErr)

	if result.result != "testhost:\noutput" {
		t.Errorf("Expected 'testhost:\\noutput', got '%s'", result.result)
	}
	if result.err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, result.err)
	}
}

func TestGetHostKeyCallbackInsecure(t *testing.T) {
	opt := common.Options{
		InsecureHost: true,
	}

	callback := getHostKeyCallback(opt)

	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that callback works (should not error in insecure mode)
	err := callback("testhost", nil, nil)
	if err != nil {
		t.Errorf("Expected nil error in insecure mode, got %v", err)
	}
}

func TestGetHostKeyCallbackDefault(t *testing.T) {
	opt := common.Options{
		InsecureHost: false,
	}

	callback := getHostKeyCallback(opt)

	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Default mode also uses InsecureIgnoreHostKey but with warning
	// The callback should still work
	err := callback("testhost", nil, nil)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestGetHostKeyCallbackWithKnownHosts(t *testing.T) {
	// Test with a known hosts path set
	opt := common.Options{
		InsecureHost: false,
		KnownHosts:   "/nonexistent/path/known_hosts",
	}

	callback := getHostKeyCallback(opt)

	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Should fall back to insecure mode since file doesn't exist
	// and return nil
	err := callback("testhost", nil, nil)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
