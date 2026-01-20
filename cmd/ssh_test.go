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

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestSSHCommand(t *testing.T) {
	// Find the ssh command
	var sshCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "ssh" {
			sshCmd = cmd
			break
		}
	}

	if sshCmd == nil {
		t.Fatal("ssh command not found")
	}

	// Test ssh command properties
	if sshCmd.Use != "ssh" {
		t.Errorf("sshCmd.Use = %s, want ssh", sshCmd.Use)
	}

	if sshCmd.Short == "" {
		t.Error("sshCmd.Short is empty")
	}

	if sshCmd.Long == "" {
		t.Error("sshCmd.Long is empty")
	}

	// Test that command flag exists
	commandFlag := sshCmd.Flags().Lookup("command")
	if commandFlag == nil {
		t.Error("command flag not found")
	}

	// Test flag shorthand
	if commandFlag.Shorthand != "c" {
		t.Errorf("command flag shorthand = %s, want c", commandFlag.Shorthand)
	}
}

func TestSSHCommandFlagBinding(t *testing.T) {
	// Test that the ssh command flag is properly bound to viper
	// Reset viper
	viper.Reset()

	// Find the ssh command
	var sshCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "ssh" {
			sshCmd = cmd
			break
		}
	}

	if sshCmd == nil {
		t.Fatal("ssh command not found")
	}

	// Check that the flag is bound to the correct viper key
	commandFlag := sshCmd.Flags().Lookup("command")
	if commandFlag == nil {
		t.Fatal("command flag not found")
	}

	// The flag should be bound to "ya.ssh.command"
	// We can't easily test the actual binding without running the command,
	// but we can verify the flag exists and has the right properties
	if commandFlag.Name != "command" {
		t.Errorf("flag name = %s, want command", commandFlag.Name)
	}
}

func TestSSHCommandUsage(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "SSH with help flag",
			args:    []string{"ssh", "--help"},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use RootCmd since it has ssh as a subcommand
			testCmd := RootCmd
			testCmd.SetArgs(tt.args)

			err := testCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSSHCommandRunFunction(t *testing.T) {
	// Test that the Run function exists and can be called
	// Find the ssh command
	var sshCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "ssh" {
			sshCmd = cmd
			break
		}
	}

	if sshCmd == nil {
		t.Fatal("ssh command not found")
	}

	// Verify Run function is set
	if sshCmd.Run == nil {
		t.Fatal("sshCmd.Run is nil")
	}

	// Test calling Run with empty args (will fail but exercises the code path)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run function panicked: %v", r)
		}
	}()

	// Set some viper values to test the Run function
	viper.Set("ya.machines", []string{"host1"})
	viper.Set("ya.ssh.command", "ls")

	// Call Run - it will fail to connect but exercises the code
	sshCmd.Run(sshCmd, []string{})
}

func TestSSHCommandStructure(t *testing.T) {
	// Find the ssh command
	var sshCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "ssh" {
			sshCmd = cmd
			break
		}
	}

	if sshCmd == nil {
		t.Fatal("ssh command not found")
	}

	// Verify the command is properly configured
	if sshCmd.Use == "" {
		t.Error("sshCmd.Use is empty")
	}

	// Verify Run is set
	if sshCmd.Run == nil {
		t.Error("sshCmd.Run is nil")
	}

	// Verify the command has the expected flag
	flag := sshCmd.Flags().Lookup("command")
	if flag == nil {
		t.Error("command flag not found")
	}
}
