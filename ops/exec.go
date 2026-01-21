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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/raravena80/ya/common"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
)

// executeResult holds the result of an SSH operation.
type executeResult struct {
	result string
	err    error
}

// Formatter defines the interface for output formatting.
type Formatter interface {
	FormatResult(hostname, output string, err error) string
	FormatError(err error) string
}

// TextFormatter implements plain text output formatting.
type TextFormatter struct{}

func (f *TextFormatter) FormatResult(hostname, output string, err error) string {
	return hostname + ":\n" + output
}

func (f *TextFormatter) FormatError(err error) string {
	return fmt.Sprintf("Error: %v", err)
}

// JSONFormatter implements JSON output formatting.
type JSONFormatter struct{}

func (f *JSONFormatter) FormatResult(hostname, output string, err error) string {
	var status string
	if err != nil {
		status = fmt.Sprintf("\"error\": %q", err.Error())
	} else {
		status = fmt.Sprintf("\"output\": %q", output)
	}
	return fmt.Sprintf(`{"host": %q, %s}`, hostname, status)
}

func (f *JSONFormatter) FormatError(err error) string {
	return fmt.Sprintf(`{"error": %q}`, err.Error())
}

// matchesPattern checks if a hostname matches a glob pattern.
// Supports wildcards: * (matches any sequence) and ? (matches any single character).
func matchesPattern(host, pattern string) bool {
	// First try filepath.Match for standard glob matching
	matched, err := filepath.Match(pattern, host)
	if matched && err == nil {
		return true
	}

	// If that doesn't work, try manual wildcard matching
	// Convert glob pattern to a simple match
	return globMatch(host, pattern)
}

// globMatch performs simple glob pattern matching without using filepath.Match.
// * matches any sequence of characters (including empty)
// ? matches any single character
func globMatch(host, pattern string) bool {
	hostIdx := 0
	patIdx := 0
	hostLen := len(host)
	patLen := len(pattern)

	// Track positions for backtracking with *
	var lastStarPatIdx, lastStarHostIdx int = -1, -1

	for hostIdx < hostLen {
		if patIdx < patLen && (pattern[patIdx] == '?' || pattern[patIdx] == host[hostIdx]) {
			// Character matches or ? wildcard
			patIdx++
			hostIdx++
		} else if patIdx < patLen && pattern[patIdx] == '*' {
			// Remember star position and current host position
			lastStarPatIdx = patIdx
			lastStarHostIdx = hostIdx
			patIdx++ // Move past the *
		} else if lastStarPatIdx != -1 {
			// We have a * to backtrack to - try matching one more character with *
			patIdx = lastStarPatIdx + 1
			lastStarHostIdx++
			hostIdx = lastStarHostIdx
		} else {
			// No match and no star to backtrack to
			return false
		}
	}

	// Handle remaining * in pattern
	for patIdx < patLen && pattern[patIdx] == '*' {
		patIdx++
	}

	return patIdx == patLen
}

// shouldIncludeHost determines if a host should be included based on patterns.
func shouldIncludeHost(host string, patterns, excludes []string) bool {
	// If no patterns specified and no excludes specified, include all hosts
	if len(patterns) == 0 && len(excludes) == 0 {
		return true
	}

	// Check exclusions first - exclusions always take priority
	for _, pattern := range excludes {
		if matchesPattern(host, pattern) {
			return false
		}
	}

	// If no inclusion patterns specified, include remaining hosts
	// (excludes have already been filtered out above)
	if len(patterns) == 0 {
		return true
	}

	// Check inclusion patterns - only include if matches a pattern
	for _, pattern := range patterns {
		if matchesPattern(host, pattern) {
			return true
		}
	}

	return false
}

// filterHosts returns the list of hosts that match the inclusion/exclusion patterns.
func filterHosts(hosts []string, patterns, excludes []string) []string {
	var filtered []string
	for _, host := range hosts {
		if shouldIncludeHost(host, patterns, excludes) {
			filtered = append(filtered, host)
		}
	}
	return filtered
}

type execFuncType func(common.Options, string, *ssh.ClientConfig) executeResult

// makeExecResult creates a new executeResult with the given hostname, output, and error.
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

// SSHSession creates SSH sessions to multiple machines and executes commands or copy operations.
// It takes functional options to configure the SSH connection and runs the operation concurrently
// on all specified machines. Returns true if all operations succeed, false otherwise.
// This function uses a background context. For cancellation support, use SSHSessionWithContext.
func SSHSession(options ...func(*common.Options)) bool {
	return SSHSessionWithContext(context.Background(), options...)
}

// SSHSessionWithContext creates SSH sessions with context support for cancellation.
// If the context is cancelled before all operations complete, the function returns early.
// Returns true if all operations succeed, false otherwise.
func SSHSessionWithContext(ctx context.Context, options ...func(*common.Options)) bool {
	var execFunc execFuncType

	opt := common.Options{}
	for _, option := range options {
		option(&opt)
	}

	// Filter hosts based on patterns
	machines := filterHosts(opt.Machines, opt.HostPatterns, opt.HostExcludes)

	// Handle dry-run mode
	if opt.DryRun {
		fmt.Fprintln(os.Stderr, "DRY-RUN: Previewing operations (no actual execution)")
		for _, m := range machines {
			if opt.Op == "ssh" {
				fmt.Printf("DRY-RUN: Would execute on %s: %s\n", m, opt.Cmd)
			} else if opt.Op == "scp" {
				if opt.IsRecursive {
					fmt.Printf("DRY-RUN: Would copy (recursive) %s to %s:%s\n", opt.Src, m, opt.Dst)
				} else {
					fmt.Printf("DRY-RUN: Would copy %s to %s:%s\n", opt.Src, m, opt.Dst)
				}
			}
		}
		return true
	}

	// Determine connection timeout
	var connectTimeout time.Duration
	if opt.ConnectTimeout != nil {
		connectTimeout = time.Duration(*opt.ConnectTimeout) * time.Second
	} else {
		connectTimeout = time.Duration(opt.Timeout) * time.Second
	}

	// done channel for synchronization
	done := make(chan bool, len(machines))

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
		Timeout:         connectTimeout,
	}

	// Get formatter based on output format
	var formatter Formatter = &TextFormatter{}
	switch opt.OutputFormat {
	case "json":
		formatter = &JSONFormatter{}
	case "yaml":
		// YAML formatter - simple implementation
		formatter = &TextFormatter{} // Fall back to text for now
	case "table":
		// Table formatter - simple implementation
		formatter = &TextFormatter{} // Fall back to text for now
	}

	for _, m := range machines {
		// we'll write results into the buffered channel of strings
		switch opt.Op {
		case "ssh":
			execFunc = executeCmd
		case "scp":
			execFunc = executeCopy
		}
		go func(hostname string, execFunc execFuncType) {
			select {
			case <-ctx.Done():
				fmt.Println(hostname, ":", ctx.Err())
				done <- false
				return
			default:
			}
			res := execFunc(opt, hostname, config)
			if res.err == nil {
				if opt.OutputFormat == "json" {
					fmt.Println(formatter.FormatResult(hostname, res.result, nil))
				} else {
					fmt.Print(res.result)
				}
				done <- true
			} else {
				fmt.Println(res.result, "\n", res.err)
				done <- false
			}
		}(m, execFunc)
	}

	retval := true
	for i := 0; i < len(machines); i++ {
		select {
		case <-ctx.Done():
			// Context was cancelled, drain remaining goroutines
			go func() {
				for j := i + 1; j < len(machines); j++ {
					<-done
				}
			}()
			return false
		case success := <-done:
			if !success {
				retval = false
			}
		}
	}
	return retval
}
