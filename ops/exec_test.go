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
	"strings"
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

func TestTextFormatter(t *testing.T) {
	f := &TextFormatter{}

	tests := []struct {
		name     string
		hostname string
		output   string
		err      error
		expected string
	}{
		{
			name:     "Format result without error",
			hostname: "testhost",
			output:   "Hello World",
			err:      nil,
			expected: "testhost:\nHello World",
		},
		{
			name:     "Format result with error",
			hostname: "testhost",
			output:   "",
			err:      fmt.Errorf("connection failed"),
			expected: "testhost:\n",
		},
		{
			name:     "Format error message",
			hostname: "testhost",
			output:   "",
			err:      fmt.Errorf("connection failed"),
			expected: "Error: connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil && tt.name == "Format error message" {
				result := f.FormatError(tt.err)
				if result != tt.expected {
					t.Errorf("FormatError() = %v, want %v", result, tt.expected)
				}
			} else {
				result := f.FormatResult(tt.hostname, tt.output, tt.err)
				if result != tt.expected {
					t.Errorf("FormatResult() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestJSONFormatter(t *testing.T) {
	f := &JSONFormatter{}

	tests := []struct {
		name     string
		hostname string
		output   string
		err      error
		contains string
	}{
		{
			name:     "Format result without error",
			hostname: "testhost",
			output:   "Hello World",
			err:      nil,
			contains: `"host": "testhost"`,
		},
		{
			name:     "Format result with error",
			hostname: "testhost",
			output:   "",
			err:      fmt.Errorf("connection failed"),
			contains: `"error": "connection failed"`,
		},
		{
			name:     "Format error message",
			hostname: "testhost",
			output:   "",
			err:      fmt.Errorf("timeout"),
			contains: `"error": "timeout"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.err != nil && tt.name == "Format error message" {
				result = f.FormatError(tt.err)
			} else {
				result = f.FormatResult(tt.hostname, tt.output, tt.err)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Format output does not contain %q, got: %s", tt.contains, result)
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		pattern  string
		expected bool
	}{
		{name: "Exact match", host: "prod-server-1", pattern: "prod-server-1", expected: true},
		{name: "Wildcard match at end", host: "prod-server-1", pattern: "prod-*", expected: true},
		{name: "Wildcard match at start", host: "prod-server-1", pattern: "*-server-1", expected: true},
		{name: "Wildcard match in middle", host: "prod-server-1", pattern: "*-*-1", expected: true},
		{name: "Wildcard match all", host: "any-host", pattern: "*", expected: true},
		{name: "Wildcard substring match", host: "prod-server-1", pattern: "prod*", expected: true},
		{name: "No match", host: "prod-server-1", pattern: "staging", expected: false},
		{name: "Partial match with wildcard", host: "prod-server-1", pattern: "prod-*", expected: true},
		{name: "Different host", host: "staging-server-1", pattern: "prod-*", expected: false},
		{name: "Question mark wildcard", host: "host-1", pattern: "host-?", expected: true},
		{name: "Multiple question marks", host: "host-12", pattern: "host-??", expected: true},
		{name: "Multiple wildcards", host: "backup-server-1", pattern: "*backup*", expected: true},
		{name: "Hyphenated wildcard pattern", host: "some-backup-server", pattern: "*-backup-*", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPattern(tt.host, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v", tt.host, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestShouldIncludeHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		patterns []string
		excludes []string
		expected bool
	}{
		{
			name:     "No patterns - include all",
			host:     "any-host",
			patterns: []string{},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "Matching pattern",
			host:     "prod-server-1",
			patterns: []string{"prod-*"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "Non-matching pattern",
			host:     "staging-server-1",
			patterns: []string{"prod-*"},
			excludes: []string{},
			expected: false,
		},
		{
			name:     "Multiple patterns, one matches",
			host:     "staging-server-1",
			patterns: []string{"prod-*", "staging-*"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "Excluded by pattern",
			host:     "prod-backup-1",
			patterns: []string{"prod-*"},
			excludes: []string{"*-backup-*"},
			expected: false,
		},
		{
			name:     "Included but excluded",
			host:     "prod-server-1",
			patterns: []string{"*"},
			excludes: []string{"*-server-*"},
			expected: false,
		},
		{
			name:     "Exclusion without inclusion",
			host:     "backup-server-1",
			patterns: []string{},
			excludes: []string{"*backup*"},
			expected: false,
		},
		{
			name:     "Multiple excludes, none match",
			host:     "prod-server-1",
			patterns: []string{"prod-*"},
			excludes: []string{"*-backup-*", "*-test-*"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIncludeHost(tt.host, tt.patterns, tt.excludes)
			if result != tt.expected {
				t.Errorf("shouldIncludeHost(%q, %v, %v) = %v, want %v",
					tt.host, tt.patterns, tt.excludes, result, tt.expected)
			}
		})
	}
}

func TestFilterHosts(t *testing.T) {
	tests := []struct {
		name      string
		hosts     []string
		patterns  []string
		excludes  []string
		expected  []string
	}{
		{
			name:     "No filters - return all",
			hosts:    []string{"host1", "host2", "host3"},
			patterns: []string{},
			excludes: []string{},
			expected: []string{"host1", "host2", "host3"},
		},
		{
			name:     "Filter by pattern",
			hosts:    []string{"prod-1", "prod-2", "staging-1"},
			patterns: []string{"prod-*"},
			excludes: []string{},
			expected: []string{"prod-1", "prod-2"},
		},
		{
			name:     "Exclude pattern",
			hosts:    []string{"prod-1", "prod-backup", "prod-2"},
			patterns: []string{},
			excludes: []string{"*-backup"},
			expected: []string{"prod-1", "prod-2"},
		},
		{
			name:     "Both include and exclude",
			hosts:    []string{"prod-1", "prod-backup", "prod-2", "staging-1"},
			patterns: []string{"prod-*"},
			excludes: []string{"*-backup"},
			expected: []string{"prod-1", "prod-2"},
		},
		{
			name:     "Multiple patterns",
			hosts:    []string{"prod-1", "staging-1", "test-1", "dev-1"},
			patterns: []string{"prod-*", "staging-*"},
			excludes: []string{},
			expected: []string{"prod-1", "staging-1"},
		},
		{
			name:     "Empty host list",
			hosts:    []string{},
			patterns: []string{"*"},
			excludes: []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterHosts(tt.hosts, tt.patterns, tt.excludes)
			if !equalStringSlices(result, tt.expected) {
				t.Errorf("filterHosts(%v, %v, %v) = %v, want %v",
					tt.hosts, tt.patterns, tt.excludes, result, tt.expected)
			}
		})
	}
}

// Helper function for string slice comparison
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
